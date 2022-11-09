package main

import "gorm.io/gen"
import "github.com/hardcaporg/hardcap/internal/db"

type Querier interface {
    // SELECT * FROM @@table WHERE id=@id
    GetByID(id int) (gen.T, error)

    // SELECT * FROM @@table WHERE sid=@sid
    GetBySID(sid string) (gen.T, error)
}

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath: "../../internal/db",
        Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface,
    })

    db.Initialize()
    g.UseDB(db.Pool)
    g.GenerateModel("registrations", gen.FieldType("id", "int64"))
    g.ApplyInterface(func(Querier){}, g.GenerateModel("registrations"))

    g.Execute()
}