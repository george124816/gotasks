
## Instalation

```sh
go install github.com/george124816/gotasks@latest
```

## Usage

```sh
# create some tasks
gotasks -action create -name "create some random project"
gotasks -action create -name "bump sql" -description "the project need to bump the latest version of sql"

# complete a task
gotasks -action complete -id 2

# show only enabled
gotasks 

# show all tasks either enable or disabled
gotasks -all

#get a specific task
gotasks -action get -id 1
```
