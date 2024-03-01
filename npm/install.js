const path = require("path")
// Just install for darwin for sake of simplicity, sorry.
// Naively include both am64/arm64 arch in the node package.

// maps process.arch to GOARCH
let GOARCH_MAP = {
  arm64: "arm64",
  amd64: "amd64",
  x86_64: "amd64",
  x64: "amd64",
  ia32: "386",
  i386: "386",
}

let GOOS_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
}

if (!(process.arch in GOARCH_MAP)) {
  console.error(`Sorry this is only packaged for ${GOARCH_MAP} at the moment.`)
  process.exit(1)
}

if (!(process.platform in GOOS_MAP)) {
  console.error(`Sorry this is only packaged for ${GOOS_MAP} at the moment.`)
  process.exit(1)
}

const arch = GOARCH_MAP[process.arch]
const platform = GOOS_MAP[process.platform]
// const installTarget = `gen-gopm3-${platform}-${arch}`
const installTarget = `envii_build_${platform}_${arch}`

const root = path.resolve(__dirname, "..")
const distDir = path.resolve(root, "dist")

const { exec } = require("child_process")
exec(`cp ./dist/${installTarget}/envii dist/envii`, (err) => {
  if (err) {
    console.error(err)
    process.exit(1)
  }
})

exec(`cp ./dist/${installTarget}/envii ../.bin/envii`, (err) => {
  if (err) {
    console.error(err)
    process.exit(1)
  }
})
