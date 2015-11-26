// suriface - Interfaces SUR con otros Sistemas Afip
package main

import (
"database/sql"
"fmt"
"os"
"io/ioutil"
"log"
"strings"
"encoding/json"
"github.com/howeyc/gopass"
_ "github.com/mattn/go-oci8"

)

type Otn struct {
   Owner string
   Type  string
   Name  string
}

type SrcDb struct {
    Db string
    Url  string
    Env  []string
    Owner  []string
    ObjType []string
    ObjName  []string
    Where  string
}

func GetConfig (cfgfile string) (Dbs []SrcDb, e error) {

file, err := ioutil.ReadFile(cfgfile)
    if err!=nil{
            log.Println("config:")
            log.Fatal(err)
    }

json.Unmarshal(file,&Dbs)
return Dbs,e
}

func armaWhere(owners []string, types []string, names []string) (w string) {
w = ""

if (len(owners) > 0) {
 w = w + " where owner in ("
   for k,v  := range owners {
    if  (k>0) {
    w =  w + ","
    }
     w = w + "'"+v+"'" 
   }
 w = w + ")"
}

if (len(types) > 0) {
if (len(w) > 0) {
 w = w + " and object_type in ("
} else {
 w = w + " where object_type in ("
}
   for k,v  := range types {
    if  (k>0) {
    w =  w + ","
    }
     w = w + "'"+v+"'" 
   }
 w = w + ")"

}

if (len(names) > 0) {
if (len(w) > 0) {
 w = w + " and object_name in ("
} else {
 w = w + " where object_name in ("
}
   for k,v  := range names {
    if  (k>0) {
    w =  w + ","
    }
     w = w + "'"+v+"'" 
   }
 w = w + ")"

}

return w
}

func genFile(d SrcDb ,o Otn, ddl string) {
fname :=  strings.ToLower(o.Owner)+"."+strings.ToLower(o.Type)+"."+strings.ToLower(o.Name)
 f, err := os.Create("./"+strings.ToLower(d.Db)+"/"+fname)
 if  err != nil { 
     log.Fatal(err)
 }
 defer f.Close()

f.WriteString(ddl)
f.Sync()
}

func DumpDb(d SrcDb) {
var o Otn
var ddl string
var askPasswd = false
var xtraSlash =""
var where = ""

fmt.Println ("Db "+d.Db)

at:=strings.IndexRune(d.Url, '@')
slash:=strings.IndexRune(d.Url, '/')

if slash < at {
  if (slash+1) == at  {
   askPasswd = true
  } 
} else {
   xtraSlash="/" 
   askPasswd = true
}

if askPasswd  {
  fmt.Print("Password: ") 
  pass:= string(gopass.GetPasswd()[:])
  urls := strings.Split(d.Url,"@")
  d.Url =urls[0]+xtraSlash+pass+"@"+urls[1]
}

// Genera directorio de salida, uno por Db
err:= os.Mkdir("./"+strings.ToLower(d.Db),0777)

// Setea entorno
for k, v := range d.Env {
var key  string

 par := k%2 == 0
 if par { 
     key=v 
 } else {
  os.Setenv(strings.ToUpper(key),strings.ToUpper(v))
 }
}

// Conecta a Oracle
  db, err := sql.Open("oci8", d.Url)
  if err != nil {
     log.Fatal(err)
      return
  }
defer db.Close()

// arma Where
if len(d.Where) > 1 {
where = d.Where
} else {
where = armaWhere(d.Owner,d.ObjType,d.ObjName)
}

// Query principal
stmt, err := db.Prepare("select dbms_metadata.get_ddl(:1,:2,:3) from DUAL")
if err != nil {
	log.Fatal(err)
}
defer stmt.Close()

// Query a DBA_OBJECTS
rows, err := db.Query("select owner,object_name,object_type from dba_objects "+ where )
if err != nil {
	log.Fatal(err)
}
defer rows.Close()
for rows.Next() {
	err := rows.Scan(&o.Owner,&o.Name,&o.Type)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(o)

        rows2, err := stmt.Query(o.Type,o.Name,o.Owner)
         if err != nil {
	    log.Fatal(err)
         }
        defer rows2.Close()
        for rows2.Next() {
	    err := rows2.Scan(&ddl)
	    if err != nil {
		log.Fatal(err)
	    }
	 genFile(d,o,ddl)
        }
        if err = rows2.Err(); err != nil {
	log.Fatal(err)
        }


}
err = rows.Err()
if err != nil {
	log.Fatal(err)
}


//    for _, v := range m {
//    }

}

func main() { 
var cfgfile string

 if len(os.Args) == 1 {
  cfgfile="./config.json"
  fmt.Println("Usando configuracion default : ./config.json")
 } else {
  if len(os.Args) < 3 {
  cfgfile=os.Args[1]
  } else {
   fmt.Fprintf(os.Stderr, "Uso: ora_dump [config.json]\n")
   os.Exit(1)
  }
 }

dbs,_:= GetConfig(cfgfile)

for _,v := range dbs {
  log.Println("Procesando Db :"+v.Db)
  DumpDb(v)
}

}
