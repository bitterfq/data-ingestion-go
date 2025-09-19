# Sep 18

1. weighted distributions for suppliers needs implementation
2. look into go tests [DONE]
3. implement parts generator [DONE]
   1. write to csv as well [DONE]
4. Improve the README [DONE]
5. Add documentation for Sep 17 work[DONE]
6. look into temporal


write to sqlite using sqlc,
add tests for that flow
if i wanted test api calls from sqlite to s3

>forget to profiling
>dont focus on getting gen perfect
>math gen later

# Sep 19 

1. clean out the db [n/a]
2. make crud app ontop of sqlite
   1. simple json api -> REST API
   2. go has http built in server
   3. insert & delete handler
   4. THEN LOOK INTO OPENAPI to replace handlers
THEN do 3
3. look into opentelemetry
4. look into openapi / swagger : https://github.com/oapi-codegen/oapi-codegen
5. sqlite api -> s3 [n/a]
6. look into temporal [n/a not needed yet]
7. improve README evenmore
8. Ask Ali about file structuring, code practices
   1. how to order, what usually to commit, the sql files and yaml files
      1. i want criticism.

side project
http call to api to get weather
   (/weather openapi endpoint) -> return to caller
   instrument it w/ opentelemetry (Trace caller, trace db, & trace api call)
   ask melvin where to send opentelemtry views


cmd server for http server, move curr main to generator in cmd dir -- std patterns