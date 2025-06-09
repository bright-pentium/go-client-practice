
TODO:
- ci 
- tests


vscode debug config:
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Echo Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "./cmd/server",
            "cwd": "${workspaceFolder}",
            "env": {},
            "args": []
        }

    ]
}
```
