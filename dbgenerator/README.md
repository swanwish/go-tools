# Simple DB Table generator

This file is the helper to generate tables from schema files.
And schema is human readable, so we can manipulate table structure in xml file, and generate table use this tool.

I am newer to golang, and I find it is powerful.
This is also an example about how to parse xml with golang.

## schema.xml
This is the xml file define the db schema.

xml fields
Table - Define the tables in the database.
Column - Define the table column, include column name, type, length primary key, and not null property.

## generator.go
This file parse the xml file, and generate the table create statement.
I comment some statements, open them can generate table on database.

The output result is like below:
```
2015/03/03 11:47:47 Execute sql:
DROP TABLE IF EXISTS Persons
2015/03/03 11:47:47 Execute sql:
CREATE TABLE Persons (
    PersonID INT NOT NULL COMMENT 'Column comment',
    LastName VARCHAR(255) NOT NULL,
    FirstName VARCHAR(255) NOT NULL,
    Address VARCHAR(255),
    City VARCHAR(255),
    PRIMARY KEY (PersonID)
) COMMENT='Test table'
```