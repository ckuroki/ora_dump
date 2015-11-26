ora_dump
========

Dump de objetos de bases de datos Oracle a Texto

Requerimientos
--------------
Necesita tener instalado un cliente oracle (puede ser instantclient)

Uso
---

El programa lee la configuracion de un archivo de configuración json. 
Por default, lo busca en el mismo directorio con el nombre `config.json`.

Si se quiere especificar otro archivo se puede pasar como parámetro.

```
ora_dump config2.json
ora_dump /usr/local/etc/mi_config.json
```

Formato del archivo de configuracion
------------------------------------

El archivo de configuracion permite definir N bases de datos.
En cada una se puede especificar una lista de Owners,Object Types o Object Names a filtrar, o directamente escribir una condicion where para la tabla DBA_OBJECTS.
Si se define el parámetro `Where`, `Owner`,`ObjType` y `ObjName` no son tomados en cuenta.

Ejemplo con 1 base, dos esquemas, todos los objetos, sin where
```
[
{
 "Db": "PROD",
 "Url": "system/manager@192.168.0.115:1521/PROD",
 "Env": ["NLS_LANG","AMERICAN_AMERICA.UTF8"],
 "Owner": ["SALES","RRHH"],
 "ObjType": [],
 "ObjName":[],
 "Where": ""
}
]
```

Ejemplo con 1 base, con condicion where
```
[
{
 "Db": "PROD",
 "Url": "system/manager@192.168.0.115:1521/PROD",
 "Env": ["NLS_LANG","AMERICAN_AMERICA.UTF8"],
 "Owner": ["SALES","RRHH"],
 "ObjType": [],
 "ObjName":[],
 "Where": " where object_type in ('FUNCTION') and owner not like 'SYS%' and last_ddl_time > sysdate-10 "
}
]
```

Ejemplo con 2 bases de datos, la segunda completa (todos los objetos)
```
[
{
 "Db": "PROD",
 "Url": "system/manager@10.2.64.115:1521/PROD",
 "Env": ["NLS_LANG","AMERICAN_AMERICA.UTF8"],
 "Owner": ["SALES","RRHH"],
 "ObjType": ["TABLE","INDEX"],
 "ObjName":[],
 "Where":""
},
{
 "Db": "test",
 "Url": "system/manager@127.0.0.1:1521/orcl",
 "Env": ["NLS_LANG","AMERICAN_AMERICA.UTF8"],
 "Owner": [],
 "ObjType": [],
 "ObjName":[],
 "Where":""
}
]
```

Funcionamiento del programa
---------------------------

1. Por cada base de datos definida en la configuración
2. Crea un directorio (si no existe con el nombre de la base)
3. Genera la lista de objetos a extraer a partir de una consulta a DBA_OBJECTS
4. Por cada objeto obtenido, utiliza DBMS_METADATA.GET_DDL para extraer la definicion 
5. Guarda el resultado en un archivo 


