package main

import (
        "fmt"
        "strconv"
        "encoding/json"
        "bytes"
        "strings"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)

// First Chaincode structure
type FirstChaincode struct{
}

// Define user structure. Structure tags are used by encoding/library
type User struct {
        Name string `json:"name"`
        SValue int `json:"svalue"` // Social
        EValue int `json:"evalue"` // Econimical
        PValue int `json:"pvalue"` // Proffesion
}

func (t *FirstChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
        fmt.Println("firstChaincode Init")

        return shim.Success(nil)
}

func (t *FirstChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Println("firstChaincode Invoke")
        // Retrieve requested function and arguments
        function, args := stub.GetFunctionAndParameters()

        if function == "createUser" {
                // Creates User
                return t.createUser(stub, args)
        } else if function == "queryUser" {
                // Query User
                return t.queryUser(stub, args)
        } else if function == "queryAllUsers" {
                // Query all Users
                return t.queryAllUsers(stub, args)
        } else if function == "sendTokens" {
                // Send Tokens
                return t.sendTokens(stub, args)
        } else if function == "deleteUser" {
                // Delete User
                return t.deleteUser(stub, args)
        }

        return shim.Error(function)
}

// '{"Args":["createUser", "USER1", "Tom", "100", "100", "100"]}'
func (t *FirstChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("firstChaincode Creating User")
        if len(args) != 5 {
                return shim.Error("Incorrect number of arguments. Expecting 5")
        }
        if (!strings.HasPrefix(args[0], "USER")) {
                fmt.Printf("- Incorrect user id")
                return shim.Error("Incorrect user id. Expecting USER* ")
        }
        var str string
        var id, s, e, p int // values
        var err error

        str = strings.TrimPrefix(args[0], "USER")
        //Check if id is integer
        id, err = strconv.Atoi(str)
        if err != nil {
                return shim.Error("Incorrect user id. Expecting integer value - USER* ")
        }
        fmt.Printf("new user id: USER%d", id)

        checkUserAsBytes, _ := stub.GetState(args[0])

        if len(checkUserAsBytes) != 0 {
                fmt.Printf("- User already exists : \n%s\n", checkUserAsBytes)
                return shim.Error("User with given key already exists")
        }

         //Check if values are integer
        s, err = strconv.Atoi(args[2])
        if err != nil {
                return shim.Error("Expecting integer value")
        }
        e, err = strconv.Atoi(args[3])
        if err != nil {
                return shim.Error("Expecting integer value")
        }
        p, err = strconv.Atoi(args[4])
        if err != nil {
                return shim.Error("Expecting integer value")
        }

        var user = User{Name: args[1], SValue: s, EValue: e, PValue: p}

        userAsBytes, _ := json.Marshal(user)
        stub.PutState(args[0], userAsBytes)
        fmt.Printf("- user created: \n%s\n", userAsBytes)

        return shim.Success(userAsBytes)
}

// '{"Args":["queryUser", "USER1"]}'
func (t *FirstChaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("firstChaincode Query User")
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

         if (!strings.HasPrefix(args[0], "USER")) {
                fmt.Printf("- Incorrect user id")
                return shim.Error("Incorrect user id. Expecting USER* ")
        }
        var str string
        var id int
        var err error

        str = strings.TrimPrefix(args[0], "USER")
        //Check if id is integer
        id, err = strconv.Atoi(str)
        if err != nil {
                return shim.Error("Incorrect user id. Expecting integer value - USER* ")
        }

        userAsBytes, _ := stub.GetState(args[0])

        if len(userAsBytes) == 0 {
                fmt.Printf("- User with given key does not exist: \n%s\n", userAsBytes)
                return shim.Error("User with given key does not exist")
        }

        fmt.Printf("- query user: USER %d \n%s\n", id, userAsBytes)
        return shim.Success(userAsBytes)
}

// '{"Args":["queryAllUsers"]}'
func (t *FirstChaincode) queryAllUsers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("firstChaincode Query All Users")

        if len(args) != 0 {
                return shim.Error("Incorrect number of arguments. Expecting 0")
        }
        startKey := "USER0"
        endKey := "USER999"

        resultsIterator, err := stub.GetStateByRange(startKey, endKey)
        if err != nil {
                return shim.Error(err.Error())
        }
        defer resultsIterator.Close()

        // buffer is a JSON array containing QueryResults
        var buffer bytes.Buffer
        buffer.WriteString("[")

        bArrayMemberAlreadyWritten := false
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()
                if err != nil {
                        return shim.Error(err.Error())
                }
                // Add a comma before array members, suppress it for the first array member
                if bArrayMemberAlreadyWritten == true {
                        buffer.WriteString(",\n")
                }
                buffer.WriteString("{\"Key\":")
                buffer.WriteString("\"")
                buffer.WriteString(queryResponse.Key)
                buffer.WriteString("\"")

                buffer.WriteString(", \"Record\":")
                // Record is a JSON object, so we write as-is
                buffer.WriteString(string(queryResponse.Value))
                buffer.WriteString("}")
                bArrayMemberAlreadyWritten = true
        }
        buffer.WriteString("]")

        fmt.Printf("- queryAllUsers: \n%s\n", buffer.String())

        return shim.Success(buffer.Bytes())
}

// '{"Args":["sendTokens", "USER1", "USER2", "100", "1"]}'

func (t *FirstChaincode) sendTokens(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("firstChaincode sending Tokens")

        if len(args) != 4 {
                return shim.Error("Incorrect number of arguments. Expecting 4")
        }

        if (!strings.HasPrefix(args[0], "USER")) {
                fmt.Printf("- Incorrect user id")
                return shim.Error("Incorrect user id. Expecting USER* ")
        }
        if (!strings.HasPrefix(args[1], "USER")) {
                fmt.Printf("- Incorrect user id")
                return shim.Error("Incorrect user id. Expecting USER* ")
        }

        var str1, str2 string
        var id1, id2 int
        var err error

        str1 = strings.TrimPrefix(args[0], "USER")
        //Check if id is integer
        id1, err = strconv.Atoi(str1)
        if err != nil {
                return shim.Error("Incorrect user id. Expecting integer value - USER* ")
        }

        str2 = strings.TrimPrefix(args[1], "USER")
        //Check if id is integer
        id2, err = strconv.Atoi(str2)
        if err != nil {
                return shim.Error("Incorrect user id. Expecting integer value - USER* ")
        }

        var val, token int // Sending value and token type
        val, err = strconv.Atoi(args[2])
        if err != nil {
                return shim.Error("Expecting integer value for token amount")
        }
        token, err = strconv.Atoi(args[3])
        if err != nil {
                return shim.Error("Expecting integer value for token type")
        }


        userAsBytes1, _ := stub.GetState(args[0])
        if len(userAsBytes1) == 0 {
                fmt.Printf("- User with given key does not exist: \n%s\n", userAsBytes1)
                return shim.Error("User with given key does not exist")
        }
        user1 := User{}

        userAsBytes2, _ := stub.GetState(args[1])
        if len(userAsBytes2) == 0 {
                fmt.Printf("- User with given key does not exist: \n%s\n", userAsBytes2)
                return shim.Error("User with given key does not exist")
        }

        user2 := User{}

        json.Unmarshal(userAsBytes1, &user1)
        json.Unmarshal(userAsBytes2, &user2)

        if token == 1 {
                if val <= user1.SValue {
                        user1.SValue = user1.SValue - val
                        user2.SValue = user2.SValue + val
                } else if val > user1.SValue {
                        return shim.Error("Invalid token value. User does not possess the provided amount of this type of tokens")
                }
        } else if token == 2 {
                if val <= user1.EValue {
                        user1.EValue = user1.EValue - val
                        user2.EValue = user2.EValue + val
                } else if val > user2.EValue {
                        return shim.Error("Invalid token value. User does not possess the provided amount of this type of tokens")
                }
        } else if token == 3 {
                if val <= user1.PValue {
                        user1.PValue = user1.PValue - val
                        user2.PValue = user2.PValue + val
                } else if val > user1.PValue {
                        return shim.Error("Invalid token value. User does not possess the provided amount of this type of tokens")
                }
        } else if token != 1 {
                return shim.Error("Invalid token number. Expecting: 1 - social, 2 - economical, 3 - proffession")
        }

        userAsBytes1, _ = json.Marshal(user1)
        userAsBytes2, _ = json.Marshal(user2)
        stub.PutState(args[0], userAsBytes1)
        stub.PutState(args[1], userAsBytes2)
        fmt.Printf("- Succesfull transaction of %d tokens TYPE - %d \n from : USER%d \n%s\n", val, token,  id1, userAsBytes1)
        fmt.Printf("- to user : USER%d \n%s\n", id2, userAsBytes2)
        return shim.Success(userAsBytes2)
}


// '{"Args":["deleteUser", "USER1"]}'
func (t *FirstChaincode) deleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        fmt.Println("firstChaincode Deleting User")
        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        if (!strings.HasPrefix(args[0], "USER")) {
                fmt.Printf("- Incorrect user id")
                return shim.Error("Incorrect user id. Expecting USER* ")
        }
        var str string
        var id int
        var err error

        str = strings.TrimPrefix(args[0], "USER")
        //Check if id is integer
        id, err = strconv.Atoi(str)
        if err != nil {
                return shim.Error("Incorrect user id. Expecting integer value - USER* ")
        }

        userAsBytes, _ := stub.GetState(args[0])

        if len(userAsBytes) == 0 {
                fmt.Printf("- User does not exist")
                return shim.Error("User with given key does not exist")
        }
        stub.DelState(args[0])
        fmt.Printf("- user deleted: USER%d", id)

        return shim.Success(userAsBytes)
}

// The main function is only relevant in unit test mode.
func main() {
        // Create a new First Chaincode
        err := shim.Start(new(FirstChaincode))
        if err != nil {
                fmt.Printf("Error creating new FirstChaincode: %s", err)
        }
}
