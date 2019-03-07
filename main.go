// deliveroo back-end sketch using go-graphql and mongo 
package main
 
import (
	// basis
	"fmt"
	"log"
	"net/http"
	// graphql
	"github.com/graphql-go/graphql"
	// mongo driver
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
	// local dependencies
	"../_scrypt"
)

// uppercase first >> public members
type UserData struct {		// only used for UserData retrieval
    USERNAME    string      // unique(?) username
	ID 			string 		// generated / public // can be bytes as well
	IP			string		// last used ip (from net.IP.String())
	EMAIL		string		// email address
	P_HASH		string 		// pwd hash
	I_HASH		string 		// acc info hash >> contains details on user
	FIRST_NAME	string		// first name
	LAST_NAME	string		// last name
	DOB_EPOCH	int64		// dob date
	REG_EPOCH	int64		// registration date
	LOG_EPOCH	int64		// last login date
}

type Asker struct {
	UD			UserData
	Waiting		bool		// is waiting for delivery ?
	OrderTime	int64		// order time
}

type Provider struct {
	UD			UserData
	Preparing	bool		// is rpeparing order ?
	OrderTime	int64		// order time
	ProcTime	int64		// time processing
}

type Courier struct {
	UD			UserData
	Riding		bool		// is riding
	OrderTime	int64		// order receipt time
	ETA			int64		// time to destination
}

// uppercase >> public members
type Order struct {
	ID 			string 		// generated / public // can be bytes as well
	*ASKER		Asker 		// asker public hash identifier >> client
	*PROVIDER	Provider	// provider public hash identifier >> restaurant
	*COURIER	Courier		// courier public hash identifier >> delivery person
	I_HASH		string		// information encoding hash
}

type AskerModel struct {
    UD    		string  `json:"ud,omitempty"`    		// filled with UserData
	Waiting		string	`json:"waiting,omitempty"`		// is waiting for delivery ?
	OrderTime	string	`json:"orderTime,omitempty"`	// order time
}

type ProviderModel struct {
    UD    		string  `json:"ud,omitempty"`    		// filled with UserData
	Preparing	string	`json:"preparing,omitempty"`	// is rpeparing order ?
	OrderTime	string	`json:"ot,omitempty"`			// order time
	ProcTime	string	`json:"proc,omitempty"`			// time processing
}

type CourierModel struct {
    UD    		string  `json:"ud,omitempty"`    		// filled with UserData
	Riding		string	`json:"riding,omitempty"`		// is riding
	OrderTime	string	`json:"ot,omitempty"`			// order receipt time
	ETA			string	`json:"eta,omitempty"`			// time to destination
}

type OrderModel struct {
	// ... replace the members pointers with hash for transcryption
}

AskerType := graphql.NewObject(graphql.ObjectConfig{
    Name: "Asker",
    Fields: graphql.Fields{
        "ud": &graphql.Field{
            Type: graphql.String,
        },
        "waiting": &graphql.Field{
            Type: graphql.String,
        },
        "orderTime": &graphql.Field{
            Type: graphql.String,
        },
    },
})

ProviderType := graphql.NewObject(graphql.ObjectConfig{
    Name: "Provider",
    Fields: graphql.Fields{
        "ud": &graphql.Field{
            Type: graphql.String,
        },
        "rpeparing": &graphql.Field{
            Type: graphql.String,
        },
        "proc": &graphql.Field{
            Type: graphql.String,
        },
        "ot": &graphql.Field{
            Type: graphql.String,
        },
    },
})

CourierType := graphql.NewObject(graphql.ObjectConfig{
    Name: "Courier",
    Fields: graphql.Fields{
        "ud": &graphql.Field{
            Type: graphql.String,
        },
        "riding": &graphql.Field{
            Type: graphql.String,
        },
        "ot": &graphql.Field{
            Type: graphql.String,
        },
        "eta": &graphql.Field{
            Type: graphql.String,
        },
    },
})

OrderType := graphql.NewObject(graphql.ObjectConfig{
    Name: "Order",
    // ...
})

// entry point
func main() {
	fmt.Println("Starting application...")
	// GraphQL data models >> crafted around the Go structs

	// insert graphql types here

	// set of queries available
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{},
	})
	// data changing mutation that can be run
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{},
	})
		schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	// graphql endpoint handles all querying and mutations >> powered by the schema
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: r.URL.Query().Get("query"),
	})
	json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":8080", nil)
}

func InitDBSession() *Session {
	_PORT := 12700
    session, err := mgo.Dial("localhost:" + strconv.Itoa(_PORT))
    if err != nil {
        panic(err)
    }
    defer session.Close()
	return session
}
