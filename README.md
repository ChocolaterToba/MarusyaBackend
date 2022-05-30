# MarusyaBackend

How do i launch this???
1. Install Docker
2. Go to server/settings/values_local.yaml, change special texts for ones you need
3. If you use dockerized db, don't change any database configs
4. $ sudo docker-compose -f docker-compose.yaml up --build -d
5. Done! Now your Marusia can connect to localhost:8080/api/marusia to get access to skills

How do i add skills???
1. Launch server per instructions above
2. Create separate files for each test, for example see server/settings/tests_example.csv
3. Send PUT Form-Data request to localhost:8080/api/test/add. Fields: quizAmount: 1, 2, 3, etc; file1; file2; file3 where fileN - files created during previous step
4. Done!