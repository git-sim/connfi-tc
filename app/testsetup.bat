REM Create Recipient Bob, logout
curl -v -XPOST -b testcookiefile.txt -c testcookiefile.txt localhost:8080/login?email=bob.smith@mail.com
curl -v -b testcookiefile.txt -c testcookiefile.txt -XPUT -d firstname="Bob" -d lastname="Smith" localhost:8080/account?email=bob.smith@mail.com
curl -v -XPOST -b testcookiefile.txt -c testcookiefile.txt localhost:8080/logout

REM Create Sender Alice
curl -v -XPOST -b testcookiefile.txt -c testcookiefile.txt localhost:8080/login?email=alice.smith@mail.com
curl -v -b testcookiefile.txt -c testcookiefile.txt -XPUT -d firstname="Alice" -d lastname="Smith" localhost:8080/account?email=alice.smith@mail.com

REM Send messages
FOR /L %%x IN (1,1,35) DO (
curl -v --cookie testcookiefile.txt -d "@testmsg.json" -H "Content-Type: application/json" -X POST http://localhost:8080/message
)

REM logout Sender
curl -v -XPOST -b testcookiefile.txt -c testcookiefile.txt localhost:8080/logout
