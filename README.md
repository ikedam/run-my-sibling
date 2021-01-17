# runsibling

## Abstract

Some windows launcher applications launch only .exe files and doesn't support launching script files like .bat files. `runsibling` works as a .exe file just to launch the same name .bat file, and helps those launcher applications launch files other than .exe files.

## Usage

1. Create your script file. We create `myscript.bat` for example:

    ```
    @ECHO OFF
    ECHO "I'm launched with %*"
    ```

2. Rename `runsibling.exe` to `myscript.exe` and put in the same directory where `myscript.bat` exists.

3. `myscript.exe` launches `myscript.bat`:

    ```
    > .\myscript.exe a b c d
    "I'm launched with a b c d"
    ```

## How it works

`runsibling.exe` finds a file with the same name to itself
and with an extension defined in the `PATHEXT` environment variable.

## Build

You can easily build `runsibling.exe` with docker-compose:

```
docker-compose run --rm build
```
