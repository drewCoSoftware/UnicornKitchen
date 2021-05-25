package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/drewCoSoftware/UnicornKitchen/settings"
	_ "github.com/lib/pq"
)

const (
	DEFAULT_PAGE_SIZE int = 5
	MAX_PAGE_SIZE     int = 25 // Prevent queries from using too many pages.
)

type PageArgs struct {
	First  int
	After  string
	Before string
	Last   int
}

func (args *PageArgs) IsBefore() bool {
	return args.Before != "" || args.Last > 0
}

func (args *PageArgs) Sanitize() error {
	args.First = clampPageSize(args.First)
	args.Last = clampPageSize(args.First)

	if args.Before != "" && args.After != "" {
		return errors.New("Ambiguous page args.  'Before' and 'After' are both set!")
	}

	return nil
}

func clampPageSize(pageSize int) int {
	if pageSize <= 0 {
		return DEFAULT_PAGE_SIZE
	}
	if pageSize > MAX_PAGE_SIZE {
		return MAX_PAGE_SIZE
	}
	return pageSize
}

// 	// NOTE: A time or event based cache would be a great idea for getting our first/last cursors for
// 	// a given table.  This would prevent extra round trips to the DB.
func GetCursorBoundaries(tableName string, idName string) (int64, int64) {
	firstQuery := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s ASC LIMIT 1", idName, tableName, idName)
	lastQuery := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s DESC LIMIT 1", idName, tableName, idName)

	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()

	var first int64
	var last int64
	firstRow := exec.QueryRow(firstQuery)
	firstRow.Scan(&first)

	lastRow := exec.QueryRow(lastQuery)
	lastRow.Scan(&last)

	return first, last
}

// This will create the database and setup the schema.
func CreateDatabase() {

	// Let's make the Database first...
	exists := dbExists(settings.DB_NAME)

	db := getPostGresDB()
	defer db.Close()

	if !exists {
		fmt.Println("The database doesn't exist, so we will create it!")

		_, err := db.Exec("CREATE DATABASE " + settings.DB_NAME)
		if err != nil {
			fmt.Println("There was an error creating the Database!")
			fmt.Println(err.Error())
		} else {
			fmt.Println("Database creation OK!")
		}
	} else {
		fmt.Println("The database already exists.")
	}

	fmt.Println("Creating the schema....")
	createSchema(true)
}

func Connect() *sql.DB {
	ops := settings.GetDatabaseOptions()

	connectionString := getConnectionString(ops)
	res := OpenConnection(connectionString)
	return res
}

func OpenConnection(connectionString string) *sql.DB {
	res, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	return res
}

func AddIngredient(ingredient *Ingredient) {
	// db := Connect()
	// defer db.Close()
	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()

	addIngredientInternal(exec, ingredient)

}

func addIngredientInternal(exec *dbExecutor, ingredient *Ingredient) {
	if !hasIngredient(exec, ingredient.Name) {
		fmt.Println("Adding the ingredient " + ingredient.Name)
		query := "INSERT INTO ingredients (Name, Description) VALUES ($1,$2)"

		_, err := exec.Exec(query, ingredient.Name, ingredient.Description)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("The ingredient '" + ingredient.Name + "' already exists!")
	}
}

// NOTE: We aren't doing anything about case insensitive search at this time.  Ideally we would
// create a case-insenstive unique index and then perform LOWER checks on name and $1.
// For the purposes of this toy applicaiton, we can live without the additional complexity.
func GetIngredient(name string) *Ingredient {
	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()

	return getIngredientInternal(exec, name)
}

func getRecipeInternal(exec *dbExecutor, name string) *Recipe {

	// NOTE: This query will need to have joins and stuff.
	// NOTE: I am sure that a proper string builder is a better way vs. concatenation which no
	// doubt creates a bunch of garbage along the way...
	query := `SELECT ri.ingredientid, i.Name, ri.amount
		      	FROM recipe_ingredients AS ri
			  	INNER JOIN recipes AS r on r.recipeid = ri.recipeid
			  	INNER JOIN ingredients as i on i.ingredientid = ri.ingredientid
				WHERE r.name = $1`

	rows, err := exec.Query(query, name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	res := &Recipe{}

	// We have ingredient data, so let's read it all in.
	for rows.Next() {

		ri := &RecipeIngredient{
			Ingredient: &Ingredient{},
		}
		rows.Scan(&ri.Ingredient.Id, &ri.Ingredient.Name, &ri.IngredientAmount)

		// NOTE: There is probably a better way to do this.  This looks like
		// it will create a lot of garbage....
		res.Ingredients = append(res.Ingredients, ri)
	}

	// Now we can get the basic recipe data that we care about...
	query = "SELECT recipeid, name, description FROM recipes WHERE Name = $1"
	row := exec.QueryRow(query, name)

	err = row.Scan(&res.Id, &res.Name, &res.Description)
	if err != nil {
		return nil
	}

	return res
}

func getIngredientInternal(exec *dbExecutor, name string) *Ingredient {

	query := "SELECT IngredientId, Name, Description FROM ingredients WHERE Name = $1"
	rows, err := exec.Query(query, name)
	defer rows.Close()

	if err != nil {
		panic(err)
	}

	// We should have one or none....
	hasMatch := rows.Next()
	if !hasMatch {
		return nil
	}

	// Deserialize the ingredient.....
	res := &Ingredient{}
	scanErr := rows.Scan(&res.Id, &res.Name, &res.Description)
	if scanErr != nil {
		panic(scanErr)
	}

	return res
}

func GetIngredients(pageArgs *PageArgs) []Ingredient {
	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()

	query := `SELECT ingredientid, name, description FROM ingredients`

	// Evaluate the page args....
	if pageArgs == nil {
		pageArgs = &PageArgs{
			First: DEFAULT_PAGE_SIZE,
		}
	}

	var queryArgs []interface{}

	// Cursor indicator....
	if pageArgs.IsBefore() {
		if pageArgs.Before != "" {
			query += " WHERE ingredientid < " + getArgIndex(queryArgs)
			queryArgs = append(queryArgs, pageArgs.Before)
		}

		query += " ORDER BY ingredientid DESC"
	} else {
		// Use an 'after query'  This is always the
		if pageArgs.After != "" {
			query += " WHERE ingredientid > " + getArgIndex(queryArgs)
			queryArgs = append(queryArgs, pageArgs.After)
		}

		query += " ORDER BY ingredientid ASC"
	}

	// Add pagination.....
	limit := pageArgs.First
	if pageArgs.Last > 0 {
		limit = pageArgs.Last
	}

	query += " LIMIT " + getArgIndex(queryArgs)
	queryArgs = append(queryArgs, fmt.Sprintf("%d", limit))

	if rows, err := exec.Query(query, queryArgs...); err == nil {

		data := make([]Ingredient, limit)
		count := 0

		for rows.Next() {
			i := &data[count]
			rows.Scan(&i.Id, &i.Name, &i.Description)

			count++
		}

		res := make([]Ingredient, count)
		offset := 0
		if pageArgs.IsBefore() {
			offset = count - 1
		}
		for i := 0; i < count; i++ {

			index := absInt(offset - i)
			res[i] = data[index]
		}

		return res

	} else {
		panic(err)
	}
}

// Because this didn't seem like a useful function according to the Go engineers.
func absInt(input int) int {
	if input < 0 {
		input = -input
	}
	return input
}

// Reverse the order of the items.
func reverse(items []Ingredient) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}

func getArgIndex(args []interface{}) string {
	res := fmt.Sprintf("$%d", len(args)+1)
	return res
}

func HasIngredient(name string) bool {
	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()
	return hasIngredient(exec, name)
}

func hasIngredient(exec *dbExecutor, name string) bool {
	return checkExists(exec, "ingredients", "Name", name)
}

func GetRecipe(name string) *Recipe {
	exec := CreateExecutor(Connect, nil)
	defer exec.Complete()

	return getRecipeInternal(exec, name)
}

// Determines if data exists on the given table where the column name/value pair matches.
func checkExists(exec *dbExecutor, tableName string, colName string, colVal string) bool {
	query := "SELECT " + colName + " FROM " + tableName + " WHERE " + colName + " = $1"
	rows, err := exec.Query(query, colVal)
	defer rows.Close()

	checkErr(err)

	res := rows.Next()
	return res
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get an ordered list of instructions for a given recipe.
func GetInstructions(recipeId int64) []string {
	dbExec := CreateExecutor(Connect, nil)
	defer dbExec.Complete()

	query := `SELECT Content, InstructionOrder FROM recipe_instructions 
			  WHERE recipeid = $1 
			  ORDER BY InstructionOrder ASC`

	if rows, err := dbExec.Query(query, recipeId); err == nil {

		res := make([]string, 0)
		var order int64
		for rows.Next() {
			var content string
			rows.Scan(&content, &order)

			res = append(res, content)
		}

		return res
	} else {
		panic(err)
	}

}

// Add a recipe to the database, saving new ingredients, etc. as we go.
func AddRecipe(recipe *Recipe) {

	fmt.Println("Adding the recipe: " + recipe.Name)

	if GetRecipe(recipe.Name) != nil {
		fmt.Println("The recipe " + recipe.Name + " already exists.  This will not be added!")
		return
	}

	// Set the appropriate refs + validate data (incomplete at this time.)
	validateRecipe(recipe)

	exec := CreateExecutor(Connect, &sql.TxOptions{ReadOnly: false})
	defer exec.Complete()

	// Add the yield ingredient as well...
	addIngredientInternal(exec, recipe.Yield.Ingredient)

	// Now the recipe entry....
	addRecipeInternal(exec, recipe)
	dbr := getRecipeInternal(exec, recipe.Name)
	if dbr == nil {
		panic(fmt.Sprintf("Could not find the recipe with name: %s", recipe.Name))
	}

	// Now all of the ingredient refs...
	for _, ri := range recipe.Ingredients {
		addIngredientInternal(exec, ri.Ingredient)
		addIngredientRef(exec, dbr.Id, ri)
	}

	for i, instr := range recipe.Instructions {
		addInstructionInternal(exec, dbr.Id, i, instr)
	}

	// If we don't flag the executor as being OK, then the transaction
	// won't get committed...
	exec.SetTransationFlag(true)

}

func addInstructionInternal(exec *dbExecutor, recipeId int64, order int, instr *RecipeInstruction) {
	query := `INSERT INTO recipe_instructions (RecipeId, InstructionOrder, Content) VALUES ($1, $2, $3)`
	_, err := exec.Exec(query, recipeId, order, instr.Content)
	if err != nil {
		panic(err)
	}
}

// Associates an ingredient, and an amount with the given recipe.
func addIngredientRef(exec *dbExecutor, recipeId int64, recipeIngredient *RecipeIngredient) {
	dbi := getIngredientInternal(exec, recipeIngredient.Ingredient.Name)
	if dbi == nil {
		panic(fmt.Sprintf("Could not find the ingredient with name: %s", recipeIngredient.Ingredient.Name))
	}

	query := "INSERT INTO recipe_ingredients (RecipeId, IngredientId, Amount) VALUES ($1, $2, $3)"
	_, err := exec.Exec(query, recipeId, dbi.Id, recipeIngredient.IngredientAmount)
	if err != nil {
		panic(err)
	}
}

func addRecipeInternal(exec *dbExecutor, recipe *Recipe) {

	// fmt.Println("Adding the recipe " + recipe.Name)
	query := "INSERT INTO recipes (Name, Description, YieldAmount, YieldIngredientId) VALUES ($1,$2,$3,$4)"

	yieldIngredient := getIngredientInternal(exec, recipe.Yield.Ingredient.Name)
	if yieldIngredient == nil {
		panic("There is no dabase entry for the yield ingredient!")
	}

	_, err := exec.Exec(query, recipe.Name, recipe.Description, recipe.Yield.Amount, yieldIngredient.Id)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("The recipe '" + recipe.Name + "' was added!")
	}

}

// NOTE: This validation is very basic....
func validateRecipe(recipe *Recipe) {
	if recipe.Yield == nil {
		panic("There is no yield ingredient listed!")
	}
	for x := range recipe.Ingredients {
		recipe.Ingredients[x].Recipe = recipe
	}
}

func query(db *sql.DB, query string) *sql.Rows {
	res, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	return res
}

func getConnectionString(ops *settings.DBOptions) string {
	fmtString := "port=%s host=%s user=%s password=%s dbname=%s sslmode=disable"
	useHost := ops.Address
	usePort := "5432"

	// Parse out the port / host from the address.
	parts := strings.Split(ops.Address, ":")
	pLen := len(parts)
	if pLen > 1 {
		useHost = parts[0]
		usePort = parts[1]
	}
	res := fmt.Sprintf(fmtString, usePort, useHost, ops.User, ops.Password, ops.Database)

	return res
}

func dbExists(dbName string) bool {
	// We want to check the postgres catalog this time around.
	db := getPostGresDB()

	defer db.Close()

	// Check for the database....
	qr, err := db.Query("SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + dbName + "'")
	if err != nil {
		panic(err)
	}
	defer qr.Close()

	hasOne := qr.Next()
	return hasOne
}

func getPostGresDB() *sql.DB {
	ops := settings.GetDatabaseOptions()
	ops.Database = "postgres"

	cs := getConnectionString(ops)
	db := OpenConnection(cs)
	return db
}

func exec(db *sql.DB, query string, confirm bool) {
	qr, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	if confirm {
		checkExec(qr)
	}
}

func checkExec(qr sql.Result) {
	count, err := qr.RowsAffected()
	if err != nil {
		panic(err)
	}
	if count < 1 {
		panic("The query did not succeed!")
	}
}
