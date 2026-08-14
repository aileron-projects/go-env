# Loading env file example

## Run cli

Run the main.go.
It shows it can take `-env` flag.

```sh
$ go run ./main.go --help
  -env string
        env file path to load (default "env.txt")
```

## Load env file

Run the main.go without any flags (use default `env.txt`).

```sh
go run ./main.go
```

The result will be

```txt
Parsed values: env.txt
---------------------------
Number of variables: 13

>> KEY: PASSWORD
>> VALUE: bar

>> KEY: SECRET_URL
>> VALUE: http://foo:bar@example.com

>> KEY: QUOTE_DOUBLE
>> VALUE: double quoted. ' can be used.

>> KEY: QUOTE_SINGLE_ESCAPE
>> VALUE: single quotation 'escaped'.

>> KEY: FOO
>> VALUE: foo

>> KEY: MULTILINE_A
>> VALUE: onetwo

>> KEY: MULTILINE_B
>> VALUE: ⬎
one
two

>> KEY: QUOTE_DOUBLE_ESCAPE
>> VALUE: double quotation "escaped".

>> KEY: URL
>> VALUE: http://example.com

>> KEY: QUOTE_SINGLE
>> VALUE: single quoted. " can be used.

>> KEY: BAR
>> VALUE: bar

>> KEY: BAZ
>> VALUE: baz

>> KEY: USERNAME
>> VALUE: foo
```
