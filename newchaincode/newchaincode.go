package main

import (
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
        "fmt"
        "bytes"
        "encoding/json"
        "strconv"
)
type Sample struct{}
type Watch struct{
    Id    string `json:"id"`
    Data  string `json:"data"`
    ConfirmFlag  string `json:"confirmflag"`
    FactoryConfirm string `json:"factoryconfirm"`
    ShopConfirm string `json:"shopconfirm"`
    SelfConfirm string `json:"selfconfirm"`
    Rewriter string `json:"rewriter"`
    Refuse string `json:"refuse"`
    Refuser string `json:"refuser"`
    Owner string `json:"owner"`
    Place string `json:"place"`
    Version int `json:"version"`
}

func (s *Sample) Init(stub shim.ChaincodeStubInterface) pb.Response{
     return shim.Success(nil)

}

func (s *Sample)  Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    fmt.Println("ex02 Invoke")
    function, args := stub.GetFunctionAndParameters()
    if function == "createWatch" {
            return s.createWatch(stub, args)
    } else if function == "queryOne" {
            return s.queryOne(stub, args)
    } else if function == "queryAll" {
            return s.queryAll(stub)
    } else if function == "changeData" {
            return s.changeData(stub, args)
    } else if function == "confirmData" {
            return s.confirmData(stub, args)
    } else if function == "getHistory" {
            return s.getHistory(stub, args)
    } else if function == "richQuery" {
            return s.richQuery(stub, args)
    } else if function == "changeDataAndChangeOwner" {
            return s.changeDataAndChangeOwner(stub, args)
    } else if function == "richQueryByOwner" {
            return s.richQueryByOwner(stub, args)
    } else if function == "changeDataAndNewOwner" {
            return s.changeDataAndNewOwner(stub, args)
    } else if function == "maintenanceData" {
            return s.maintenanceData(stub, args)
    } else if function == "c3RichQuery" {
            return s.c3RichQuery(stub, args)
    }else if function == "checkConfirmFlag" {
            return s.checkConfirmFlag(stub, args)
    }


    return shim.Error("Invalid invoke function name. Expecting \"queryOne\" \"queryAll\" \"changeData\" \"confirmData\" \"getHistory\"  \"richQuery\" \"changeDataAndChangeOwner\" \"richQueryByOwner\"  \"changeDataAndNewOwner\" " )
}


func (s *Sample) createWatch(stub shim.ChaincodeStubInterface, args[]string) pb.Response{

    if len(args) != 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2!")
    }

    //check
    WatchAsByte, err:= stub.GetState(args[0])
    if err == nil && WatchAsByte != nil {
        return shim.Error("Data is existed!")
    }
    // Initialize the chaincode
    A := Watch {   // Entities
        Id: args[0],
        Data: args[1],
        ConfirmFlag: "0",
        FactoryConfirm: "",
        ShopConfirm: "",
        SelfConfirm: "",
        Rewriter: "",
        Refuse: "0",
        Refuser:"",
        Owner: "",
        Place:"c2",
        Version: 0,
    }
    jsons, err := json.Marshal(A)

    // Write the state to the ledger
    err = stub.PutState(args[0], jsons)
    if err != nil {
            return shim.Error(err.Error())
    }
    return shim.Success(nil)
}

func (t *Sample) queryOne(stub shim.ChaincodeStubInterface , args []string) pb.Response {
    json, _:= stub.GetState(args[0])
    return shim.Success(json)
}

func (t *Sample) queryAll(stub shim.ChaincodeStubInterface) pb.Response {
    startKey := ""
    endKey :="~~~~~~~~~~~~~"
    rs, err := stub.GetStateByRange(startKey,endKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer rs.Close()
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    for rs.HasNext() {
        queryRS, err :=rs.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"key\":")
        buffer.WriteString("\"")
        buffer.WriteString(queryRS.Key)
        buffer.WriteString("\"")
        buffer.WriteString(", \"record\": ")
        buffer.WriteString(string(queryRS.Value))
        buffer.WriteString("}")
    }
    buffer.WriteString("]")
    return shim.Success(buffer.Bytes())
}

func (s *Sample) changeData(stub shim.ChaincodeStubInterface, args []string) pb.Response{

   if len(args) != 6 {
       return shim.Error("Incorrect number of arguments. Expecting 6! ")
   }
   Version, err := strconv.Atoi(args[4])
   if err != nil {
        return shim.Error("Expecting integer value for asset holding! ")
   }
   WatchAsByte, _:= stub.GetState(args[0])
   watch := Watch{}
   json.Unmarshal(WatchAsByte, &watch)
   
   if Version != watch.Version{
       return shim.Error("Optimistic Lock is locked! Please update the info. " )
   }

   watch.Data = args[1]
   watch.ConfirmFlag = "1"
   function := args[2]
   
   place := args[5]
   if place == "c2" {
        if function == "c1"{
             watch.FactoryConfirm = args[3]
             watch.ShopConfirm = ""
             watch.SelfConfirm = "dafault"
             watch.Rewriter = "c1"
        } else if function == "c2" {
             watch.ShopConfirm = args[3]
             watch.FactoryConfirm = ""
             watch.SelfConfirm = "dafault"
             watch.Rewriter = "c2"
        } else {
             return shim.Error("Invalid charactor name. Expecting \"c1\" \"c2\" \"c3\"." )
        }
   } else {
        if function == "c1"{
             watch.FactoryConfirm = args[3]
             watch.ShopConfirm = ""
             watch.SelfConfirm = ""
             watch.Rewriter = "c1"
        } else if function == "c2" {
             watch.ShopConfirm = args[3]
             watch.FactoryConfirm = ""
             watch.SelfConfirm = ""
             watch.Rewriter = "c2"
        } else if function == "c3" {
             watch.SelfConfirm = args[3]
             watch.FactoryConfirm = ""
             watch.ShopConfirm = ""
             watch.Rewriter = "c3"
        } else {
             return shim.Error("Invalid charactor name. Expecting \"c1\" \"c2\" \"c3\"." )
        }
   }
   watch.Refuse="0"
   watch.Version = watch.Version + 1
   WatchAsByte, _ = json.Marshal(watch)
   stub.PutState(args[0], WatchAsByte)
   return shim.Success(WatchAsByte)
}

func (s *Sample) changeDataAndChangeOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response{

   if len(args) != 6 {
       return shim.Error("Incorrect number of arguments. Expecting 6! ")
   }
   Owner := args[5]
   Version, err := strconv.Atoi(args[4])
   if err != nil {
        return shim.Error("Expecting integer value for asset holding")
   }
   WatchAsByte, _:= stub.GetState(args[0])
   watch := Watch{}
   json.Unmarshal(WatchAsByte, &watch)
   
   if Version != watch.Version{
       return shim.Error("Optimistic Lock is locked! Please update the info. " )
   }

   watch.Data = args[1]
   watch.ConfirmFlag = "1"


   watch.ShopConfirm = args[2]
   watch.FactoryConfirm = ""
   
   watch.Rewriter = "c2"

   if Owner=="" {
       watch.Place= "c2"
       watch.SelfConfirm = "default"
   }else {
       watch.Place= "c3"
       watch.SelfConfirm = ""
   }
   watch.Refuse="0"
   watch.Owner = Owner
   watch.Version = watch.Version + 1
   WatchAsByte, _ = json.Marshal(watch)
   stub.PutState(args[0], WatchAsByte)
   return shim.Success(WatchAsByte)
}

func (s *Sample) changeDataAndNewOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response{

   if len(args) != 5 {
       return shim.Error("Incorrect number of arguments. Expecting 5! ")
   }
   Owner := args[3]
   Version, err := strconv.Atoi(args[2])
   if err != nil {
        return shim.Error("Expecting integer value for asset holding! ")
   }
   WatchAsByte, _:= stub.GetState(args[0])
   watch := Watch{}
   json.Unmarshal(WatchAsByte, &watch)
   
   if Version != watch.Version{
       return shim.Error("Optimistic Lock is locked! Please update the info. " )
   }

   watch.Data = args[1]
   watch.ConfirmFlag = "0"
   if Owner=="" {
       watch.Place= "c2"
   }else {
       watch.Place= "c3"
   }
   
   watch.Owner = Owner
   watch.Version = watch.Version + 1
   WatchAsByte, _ = json.Marshal(watch)
   stub.PutState(args[0], WatchAsByte)
   return shim.Success(WatchAsByte)
}

func (s *Sample) confirmData(stub shim.ChaincodeStubInterface, args []string) pb.Response{

   if len(args) != 5 {
       return shim.Error("Incorrect number of arguments. Expecting 5. ")
   }

   Version, err := strconv.Atoi(args[3])
   if err != nil {
        return shim.Error("Expecting integer value for asset holding! ")
   }
   WatchAsByte, _:= stub.GetState(args[0])
   watch := Watch{}
   json.Unmarshal(WatchAsByte, &watch)
   if Version != watch.Version{
       return shim.Error("Optimistic Lock is locked! Please update the info. " )
   }
   refuseflag := args[4]
   if refuseflag == "1"{
        watch.Refuse = "1"
        watch.Refuser = args[2]
        watch.Version = watch.Version + 1
        WatchAsByte, _ = json.Marshal(watch)
        stub.PutState(args[0], WatchAsByte)
        return shim.Success(WatchAsByte)
   }
   function := args[1]
   if function == "c1"{
        watch.FactoryConfirm = args[2]
   } else if function == "c2" {
        watch.ShopConfirm = args[2]
   } else if function == "c3" {
        watch.SelfConfirm = args[2]
   } else {
        return shim.Error("Invalid charactor name. Expecting \"c1\" \"c2\" \"c3\"! " )
   }
   watch.Version = watch.Version + 1
   WatchAsByte, _ = json.Marshal(watch)
   stub.PutState(args[0], WatchAsByte)
   return shim.Success(WatchAsByte)
}

func (s *Sample) checkConfirmFlag(stub shim.ChaincodeStubInterface, args []string) pb.Response{
   if len(args) != 1 {
       return shim.Error("Incorrect number of arguments. Expecting 1. ")
   }
   WatchAsByte, _:= stub.GetState(args[0])
   watch := Watch{}
   json.Unmarshal(WatchAsByte, &watch)
   if ( watch.FactoryConfirm !="" && watch.ShopConfirm !="" && watch.SelfConfirm != "" ){
        watch.ConfirmFlag = "0"
        watch.FactoryConfirm =""
        watch.ShopConfirm = ""
        watch.SelfConfirm =""
        fmt.Printf("Flag check pass!! Flag = 0 Confirm = \"\" ")
        WatchAsByte, _ = json.Marshal(watch)
        stub.PutState(args[0], WatchAsByte)
   }
   return shim.Success(WatchAsByte)
}

func (t *Sample) getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    keyIter, err := stub.GetHistoryForKey(args[0])
    if err != nil {
        return shim.Error("history err")
    }

    defer keyIter.Close()
    var keys []string
    for keyIter.HasNext() {
        response, err := keyIter.Next()
        if err != nil {
                  return shim.Error("history app err")
        }
        txid := response.TxId
        txValue := response.Value
        txStatus := response.IsDelete
        fmt.Printf("id : %s , value : %s , isdelete : %s \n " , txid ,string(txValue) , txStatus )
        keys = append( keys , string(txValue) )
    }
    jsonkeys, err := json.Marshal(keys)
    if err != nil {
        return shim.Error("history jsonkeys err")
    }
    return shim.Success(jsonkeys)
}

func (t *Sample) richQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 4 {
       return shim.Error("Incorrect number of arguments. Expecting 4. ")
    }
    flag := args[0] 
    regex := args[1]
    rewriter := args[2]
    whoconfirm := args[3]

    var queryString string

    queryString = fmt.Sprintf(`{"selector":{"_id":{"$regex": "%s"},"confirmflag":"%s","$or": [{"refuse":"0","%s":""},{"refuse": "1","rewriter":"%s"}]}}` , regex , flag, whoconfirm, rewriter)
    
    rs, err := stub.GetQueryResult(queryString)

    if err != nil {
         return shim.Error("Rich query failed")
    }
    defer rs.Close()
    var buffer bytes.Buffer


    bArrayMemberAlreadyWritten := false
    buffer.WriteString(`{"result":[`)

    for rs.HasNext(){
         queryResponse, err :=rs.Next()
         if err != nil {
               return shim.Error("Fail")
         }
         if bArrayMemberAlreadyWritten  == true {
               buffer.WriteString(`,`)
         }
         buffer.WriteString(string(queryResponse.Value))
         bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString(`]}`)
    fmt.Print(`{"queryRs: %s"}}`, buffer.String())
    return shim.Success(buffer.Bytes())
}

func (t *Sample) c3RichQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 2 {
       return shim.Error("Incorrect number of arguments. Expecting 2. ")
    }
    id := args[0]
    owner := args[1]

    queryString := fmt.Sprintf(`{"selector":{"_id":{"$regex": "%s"},"confirmflag":"1", "owner":"%s","$not": {"selfconfirm":"%s"},"refuse":"0" }}`, id, owner, owner)
    
    rs, err := stub.GetQueryResult(queryString)
    if err != nil {
         return shim.Error("Rich query failed")
    }
    defer rs.Close()
    var buffer bytes.Buffer

    bArrayMemberAlreadyWritten := false
    buffer.WriteString(`{"result":[`)
    for rs.HasNext(){
         queryResponse, err :=rs.Next()
         if err != nil {
               return shim.Error("Fail")
         }
         if bArrayMemberAlreadyWritten  == true {
               buffer.WriteString(`,`)
         }
         buffer.WriteString(string(queryResponse.Value))
         bArrayMemberAlreadyWritten = true
    }

    buffer.WriteString(`]}`)

    fmt.Print(`{"queryRs: %s"}}`, buffer.String())
    return shim.Success(buffer.Bytes())
}

func (t *Sample) richQueryByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 1 {
       return shim.Error("Incorrect number of arguments. Expecting 1. ")
    }

    owner := args[0]
    queryString := fmt.Sprintf(`{"selector":{"owner": "%s"}}` , owner)
    rs, err := stub.GetQueryResult(queryString)

    if err != nil {
         return shim.Error("Rich query failed! ")
    }
    defer rs.Close()
    var buffer bytes.Buffer


    bArrayMemberAlreadyWritten := false
    buffer.WriteString(`{"result":[`)

    for rs.HasNext(){
         queryResponse, err :=rs.Next()
         if err != nil {
               return shim.Error("Fail")
         }
         if bArrayMemberAlreadyWritten  == true {
               buffer.WriteString(`,`)
         }
         buffer.WriteString(string(queryResponse.Value))
         bArrayMemberAlreadyWritten = true
    }

    buffer.WriteString(`]}`)


    fmt.Print(`{"queryRs: %s"}}`, buffer.String())
    return shim.Success(buffer.Bytes())

}
func (s *Sample) maintenanceData(stub shim.ChaincodeStubInterface, args []string) pb.Response{

   if len(args) != 4 {
       return shim.Error("Incorrect number of arguments. Expecting 4! ")
   }
   version, err := strconv.Atoi(args[3])
   if err != nil {
        return shim.Error("Expecting integer value for asset holding! ")
   }
   WatchAsByte, _:= stub.GetState(args[0])
   watch := Watch{}
   json.Unmarshal(WatchAsByte, &watch)
   
   if version != watch.Version{
       return shim.Error("Optimistic Lock is locked! Please update the info. " )
   }
   watch.Owner= args[2]
   watch.Data = args[1]
   watch.Version = watch.Version + 1
   WatchAsByte, _ = json.Marshal(watch)
   stub.PutState(args[0], WatchAsByte)
   return shim.Success(WatchAsByte)
}

func main() {
    err := shim.Start(new(Sample))
    if err != nil {
        fmt.Printf("main err : %s",err)
    }
}


