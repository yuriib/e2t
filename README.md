**e2t (entity-2-table)** is a simple code generator which is used in conjunction with ORM 
and is responsible for:
* generating entity metadata objects that simplify constructing of the SQL queries
* TBD

**NOTE**: it's a playground project for learning code generation 

**PS**: If you're, for some reason, interested in idea that was used
in the project, please, feel free to use it in your own projects,
but do not forget to leave a comment or like :)

**How to generate entity metadata objects that could be used in ORM criteria builder?**

1. In file that contains entity mapping add following command
    ```
    //go:generate e2t -entity=$GOFILE -table=YOUR_TABLE_NAME
    ```

2. Add tags to the entity fields that you want to map to the table columns
    * ``entity:"COLUMN_NAME"`` - map table column to the field. I.e.: `entity:"ID"`
    * ``join-entity:"AnotherEntity"`` - used as reference to another entity. I.e.: `join-entity:"Address"`

3. Run next command
    ```
    go generate ./path_to_your/entity_package
    ```

After that you'll get **q_\*.go** files that will contain mappings to the 
table

**CAUTION!** Generator is not as smart as you might think and cannot build entity 
metadata objects for multiple ``struct``'s that are stored in single file. 
Please split your entities into multiple files instead

**Usage of entity metadata objects**
Lets assume we have some criteria builder that accepts varargs for select
statement

I.e.: 
```
c.Select(
    "ID",
    "LAST_NAME",
    "FIRST_NAME",
).
From("USERS").
Find()
```

It looks not good at first glance, because you coupling your query to 
the table schema that could be changed and after schema changes you need to 
update all the queries that use explicit reference to table columns

Here the code how it can look if use generated entity metadata objects:

```
Select(
    QUser.Id,
    QUser.LastName,
    QUser.FirstName,
).
From(QUser.TableName()).
Find()
```

This example uses least error pron approach as you no more tight to the table schema

Now if your table schema was updated the only thing that you need to do is just modify 
your entity mapping and re-run generator

For additional info please check ``samples`` folder