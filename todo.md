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

1. make crud app ontop of sqlite [DONE]
   1. simple json api -> REST API
   2. go has http built in server
   3. insert & delete handler
   4. THEN LOOK INTO OPENAPI to replace handlers

side project
http call to api to get weather
   (/weather openapi endpoint) -> return to caller
   instrument it w/ opentelemetry (Trace caller, trace db, & trace api call)
   ask melvin where to send opentelemtry views

# Sep 22

1. look into openapi [LOW-PRIORITY]
2. look into openapi / swagger [LOW-PRIORITY]
3. build out api to get weather, store it in a json format right now
4. Improve README
5. look into NLP to SQL validation systems/frameworks
   1. Assume models are given enough context
   2. GOAL: test the output from
   3. high level document to priotize workflow
      1. describing test case how would we do it
      2. assume context is already setup , db etc and models and there is expected data and compare with return data
      3. how do u describe it
      4. test specificaion for evaluation nlp to sql
      5. how do u specify a test -> sentence.txt, query.txt, data.txt and expected result. each dir is its own test case
      6. "WHAT CAN WE DO IN ONE DAY"
6. Look into melvins llmjury project how we can link the two
7. go connector for snowflake -- Melvin -- go lang cli for snowflake [high-priority] --> temporal consistent runner

# Sep 24

end to end code , http server , fetch data from snowflake and insert into database and respond w/ data that was fetched --> real reason is to do the tracing, opentelemetry

add endpoint to server, that queries snowflake takes a query param returns data associated with it, and insert it into the database. trace the entire thing

look into structured json logging