# Antd design pro

## Upgrade

- install go by gvm

  ```
  zsh < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  gvm install go1.10beta1 -B
  gvm use go1.10beta1 --default
  ```

- nvm global package

  ```
  NPM_PACKAGES="${HOME}/.npm-packages"
  PATH="$NPM_PACKAGES/bin:$PATH"
  ```

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
  BROWSER=none PORT=3000 npm start # start
  ```
