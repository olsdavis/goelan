**Just log**
=
Since I have not found a single logging library for Go that works correctly – and I assure you that I have tried a lot of them –, I had to make my own one. This library _just works_. I made it quite quickly in order to be able to use it directly, so if you want to change something, you can make a pull request. You have everything you need in it; although, if you want a new feature, you can fork it or make a feature request as an issue.

How to use the library?
-
First, install it: 
```bash
go get github.com/olsdavis/justlog
```
Then, create your logger with the parameters you want:
```go
import "github.com/creart/justlog"

// if you want to create an empty logger
logger := justlog.New()

// if you want to create a logger which prints the messages to the console
logger := justlog.NewWithHandlers(justlog.NewConsoleHandler())

// you can set formats to the logger
logger := justlog.NewWithHandlers(justlog.NewConsoleHandler()).
        SetFormatters(NewFormatter("%{LEVEL}: %{MESSAGE}")) 
logger.Debug("Testing Creart's library!")
        // this will print: DEBUG: Testing Creart's library!
// see handler.go to see which arguments are currently available for the formatting
```
Finally, you can create your own handlers and formatters: you just have to create
a struct which implements the needed interface (Handler or Formatter).
Currently, you have the following implementations:
1. ConsoleHandler, prints the message to the console;
2. FileHandler, prints the message to the given file. (Have a look at the handler.go file.)