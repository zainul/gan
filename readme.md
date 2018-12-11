# Gan

[![Go Report Card](https://goreportcard.com/badge/github.com/zainul/gan)](https://goreportcard.com/report/github.com/zainul/gan)

Gan is migration tool and seeder tool , currently we just supported only postgres

## Installation

For installation 

```
go get -u github.com/zainul/gan
```

And the make the binary CLI with go to `$GOPATH/src/github.com/zainul/gan` and run

```
./install.sh
```

Test the installation with type
```
gan
```

will be output some help info in your command line

```
NAME:
   gan - gan use for migrate and seed the database

USAGE:
   gan [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     migrate                     Migrate migrations script
     seed                        Seed the data from file
     create_seed_file            Create seed template file
     create_migration_from_file  Create migration from SQL file
     create_migration            Create migration file
     help, h                     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE
   --help, -h              show help
   --version, -v           print the version
```

## Usage

first rule is , the config file is a must passing with argument `c`. sample of  file `config.json`

```
{
    "dir":"/home/your_pc/work/app/workspace/src/github.com/zainul/gan/examples",
    "conn":"user=postgres password=*** dbname=dbname host=localhost port=5432 sslmode=disable",
    "seed_dir":"/home/your_pc/work/app/workspace/src/github.com/zainul/gan/examples/seed"
}
```

### migrate

    migrate use for migrate all pending migration to the db, the command will be update the migration tables also

    example : gan -c "example/config.json" migrate

### seed

    seed use for seeder data to the database, with config `seed_dir` in `config.json` seeder can use from file like .json file 

    example : gan -c "example/config.json" seed

### create_seed_file

    create_seed_file use for create seeder file template to define the structure of object that will be use to seed

    example : gan -c "example/config.json" create_seed_file yourtablename

### create_migration_from_file

    create_migration_from_file is use for create migration template but the migration schema is from SQL file.

    example : gan -c "example/config.json" create_migration_from_file yourtablename

### create_migration

    create_migration_from_file is use for create migration template

    example : gan -c "example/config.json" create_migration yourtablename
