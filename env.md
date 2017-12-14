# For Development

- install go by gvm

  ```
  zsh < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  gvm install go1.4.3 -B
  gvm use go1.4.3
  export GOROOT_BOOTSTRAP=$GOROOT
  gvm install go1.10beta1
  gvm use go1.10beta1 --default
  ```

- nvm global package path

  ```
  NPM_PACKAGES="${HOME}/.npm-packages"
  PATH="$NPM_PACKAGES/bin:$PATH"
  ```

## for ant design

- change backend server(add to file .env.development.local)

  ```
  REACT_APP_BACKEND="http://localhost:8080"
  ```

## for ant design pro

- change backend server

  ```
  export default {'GET /api/*': 'http://localhost:8080/api/'};
  ```

- fix git precommit, remove this line from file package.json

  ```
  "precommit": "npm run lint-staged",
  ```

- Start without open browser and in listen at port 3000

  ```
  BROWSER=none PORT=3000 npm start
  ```
