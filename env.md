# For Development

- install go by gvm

  ```
  zsh < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  gvm install go1.4.3 -B
  gvm use go1.4.3
  export GOROOT_BOOTSTRAP=$GOROOT
  gvm install go1.10beta2
  gvm use go1.10beta2 --default
  ```

- nvm global package path

  ```
  NPM_PACKAGES="${HOME}/.npm-packages"
  PATH="$NPM_PACKAGES/bin:$PATH"
  ```

- install atom editor

  ```
  sudo pacman -S --needed gconf base-devel git nodejs npm libsecret python2 libx11 libxkbfile
  git clone --depth=1 https://github.com/atom/atom.git
  cd atom
  script/bootstrap
  script/build --compress-artifacts
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
