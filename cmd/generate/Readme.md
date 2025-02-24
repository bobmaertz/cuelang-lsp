

## Build Log 

Dec 19: I had an epiphany today - self referential structures can exist if a pointer is used. A pointer is a different type and the compiler can handle that. Feel really silly to learn this so late but very handy! 
    - SelectionRange is optional but i dont have support for optionals yet. The optional is handled in gopls via a pointer and thats how they were able to bypass the issue i've been banging my head on


## TODO: 
- [x] Cleanup main.go and add tests for existing funcs
- Add suppport for Notifications
- Add support for Requests / Responses
- 