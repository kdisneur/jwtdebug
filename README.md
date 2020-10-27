# jwtdebug

1. After cloning this project, run:

   ```
   grep -rl 'jwtdebug' . | xargs sed -i '' -e 's/jwtdebug/YOUR GITHUB REPOSITORY/g'
   ```

2. Create a new Secret for GitHub action named `MACHINE_USER` which contains a token to commit on your personal Hombrew repository
