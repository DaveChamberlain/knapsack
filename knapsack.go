package main

import "fmt"
import "os"
import "strconv"

var debug bool = false

// A structure to describe the items
type Item struct {
    weight int
    name string
    value int
}

// The information about about what is in the cell of the matrix, the current value,
// list of items, weight, and where did it come from
type Cell struct {
    val int
    items []Item
    weight int
}
func (self Cell) init() {
    self.val = 0
    self.weight = 0
}

// The matrix, it is a 2 dimensional array that stores the best combination of items
// "so far".  
type Matrix struct {
    matrix [][]Cell
    item_list []Item
    max_weight int
}
func (self *Matrix) init(item_names []Item, max_weight_allowed int) {
    var item_index int
    var i int

    // we're building the matrix as we go, we want it to be as tall as the
    // number of items (consider this the 'y' axis
    self.matrix = make([][]Cell, len(item_names))

    self.item_list = item_names
    self.max_weight = max_weight_allowed

    // We want it to be as wide as the possible weights, which will be 0 through
    // max weight.  Each cell will contain the items that make up that weight and the
    // value of those items.  It will always be the best combination of items known
    // "so far" as we go through each item.  However, here, we're just initializing it.
    for item_index = 0; item_index < len(self.item_list); item_index++ {
        self.matrix[item_index] = make([]Cell, self.max_weight + 1)
        for i = 0; i < self.max_weight + 1; i++ {
            self.matrix[item_index][i] = Cell{}
            self.matrix[item_index][i].init()
        }
    }
}
func (self *Matrix) num_items() int {
    // getter: number of items
    return len(self.item_list)
}
func (self *Matrix) get_cell(row int, col int) Cell {
    // getter: the information for a particular cell
    return self.matrix[row][col]
}
func (self *Matrix) set_cell(row int, col int, cell Cell) {
    // setter, store information about a cell
    self.matrix[row][col] = cell
}
func (self *Matrix) set_cell_weight(row int, col int, new_weight int) {
    // setter: set the weight for a particular cell in the matrix
    self.matrix[row][col].weight = new_weight
}
func (self *Matrix) set_cell_val(row int, col int, val int) {
    // setter: set the value for a particular cell in the matrix
    self.matrix[row][col].val = val
}
func (self *Matrix) set_cell_items(row int, col int, new_item []Item) {
    // setter: set the contents of a cell in the matrix
    self.matrix[row][col].items = new_item
}
func (self *Matrix) dump() {
    // accessor: Dump the contents of the knapsack on the floor and examine them
    var i int
    var w int = 0

    fmt.Println("The knapsack can hold a max weight of: ", self.max_weight)
    fmt.Println("After filling the knapsack, we have placed the following items in it: ")

    contents := self.get_cell(self.num_items() - 1, self.max_weight).items
    for i = 0; i < len(contents); i++ {
        fmt.Printf("   Name: %8s  Weight: %3d   Value: %3d\n", contents[i].name, contents[i].weight, contents[i].value)
        w += contents[i].weight
    }
    fmt.Println("This gives a total weight of the knapsack: ", w)
    fmt.Println("In the end, the value of the goods in the knapsack is: ", self.get_cell(self.num_items() - 1, self.max_weight).val)
}

// The backpack/knapsack struct keeps track of the matrix that manages the backpack, the possible items, and the max weight
type Backpack struct {
    max_weight int
    item_list []Item
    matrix Matrix
}
func (self *Backpack) init(item_list []Item, max_weight_allowed int) {
    var i int

    self.max_weight = max_weight_allowed + 1

    // Set the item list so that it is 1-based instead of 0-based, do this by
    // inserting an empty items into the list.  This makes it easier to go through
    // the list of items knowing that the top row and left column are all empty.
    self.item_list = make([]Item, len(item_list) + 1)

    for i = len(item_list); i > 0; i-- {
        self.item_list[i] = item_list[i - 1]
    }
    self.item_list[0] = Item{0, "null", 0}

    // Create the Matrix and initialize it.
    self.matrix = Matrix{}
    self.matrix.init(self.item_list, max_weight_allowed)

    // Now go ahead and fill the backpack
    self.fill_backpack()
}
func (self *Backpack) dump() {
    // accessor: dump the contents
    self.matrix.dump()
}
func (self *Backpack) fill_backpack() {
    var item int
    var possible_weight int
    var item_weight int
    var item_value int

    if debug {fmt.Println("Filling the backpack")}

    // Go through each item
    for item = 1; item < len(self.item_list); item++ {   
        if debug {fmt.Println("looking at ", item)}
        self.matrix.set_cell(item, 0, Cell{val:0})
        // Go through each possible weight of a filled backpack
        for possible_weight = 1; possible_weight < self.max_weight; possible_weight++ {
            if debug {fmt.Println("   possible_weight ", possible_weight)}

            // For convenience, get some information about the current item
            item_weight = self.item_list[item].weight
            item_value = self.item_list[item].value

            // If the weight of the current item is less than the possible weight than this item CAN go in
            // the backpack...but should it?
            if item_weight <= possible_weight {
                if debug {fmt.Println("      looks like it will fit in the backpack")}
                
                // If the value of this item will increase the value that's in the backpack for items of this weight,
                // then we're going to add this item into the backpack.
                if item_value + self.matrix.get_cell(item - 1, possible_weight - item_weight).val > self.matrix.get_cell(item - 1, possible_weight).val {
                    if debug {fmt.Println("         and it will increase the value")}
                    self.matrix.set_cell_weight(item, possible_weight, item_weight)
                    self.matrix.set_cell_items(item, possible_weight, append(self.matrix.get_cell(item - 1, possible_weight - item_weight).items, self.item_list[item]))
                    self.matrix.set_cell_val(item, possible_weight, item_value + self.matrix.get_cell(item - 1, possible_weight - item_weight).val)
                    if debug {fmt.Println("                  knapsack now contains", self.matrix.get_cell(item, possible_weight))}
                } else {
                    // If this item won't increase the value, then it doesn't need to be added, so use the best known
                    // combination "so far" by taking the information from the cell directly above this one (last item, same weight)
                    if debug {fmt.Println("         but it will not increase the value")}
                    self.matrix.set_cell_items(item, possible_weight, self.matrix.get_cell(item - 1, possible_weight).items)
                    self.matrix.set_cell_val(item, possible_weight, self.matrix.get_cell(item - 1, possible_weight).val)
                }
            } else {
                // Lastly, if the item is too heavy, then use the known information about this weight from the cell above
                self.matrix.set_cell_items(item, possible_weight, self.matrix.get_cell(item - 1, possible_weight).items)
                self.matrix.set_cell_val(item, possible_weight, self.matrix.get_cell(item - 1, possible_weight).val)
            }
        }
    }
}

func help() {
    fmt.Println("   ", os.Args[0], " [-h] [-d] <weight>")
    fmt.Println("      -h       Print this help message")
    fmt.Println("      -d       Turn on debug mode to see a narrative of the process")
    fmt.Println("")
    fmt.Println("      weight   The max weight in the knapsack")
}

//
// The main function, which has the objects, the weight, and the value.
func main() {
    var items = [6]Item{{weight: 2, name: "Item1", value: 5}, 
                        {weight: 3, name: "Item2a", value: 3}, 
                        {weight: 3, name: "Item2", value: 2}, 
                        {weight: 7, name: "Item3", value: 12}, 
                        {weight: 7, name: "Item4", value: 18}, 
                        {weight: 8, name: "Item5", value: 20}}
    var i int
    var max_weight int64 = 0

    if len(os.Args) < 2 {
        fmt.Println("Please provide the max weight of the knapsack")
        help()
        os.Exit(1)
    }
    for i = 1; i < len(os.Args); i++ {
        if os.Args[i] == "-h" {
            help()
            os.Exit(0)
        } else if os.Args[i] == "-d" {
            debug = true
        }
    }

    max_weight, _ = strconv.ParseInt(os.Args[len(os.Args)-1], 10, 0)
    b := Backpack{}
    b.init(items[:], int(max_weight))
    b.dump()
}
