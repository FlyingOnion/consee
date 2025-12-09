# Consee

Consee is a Consul management tool. It provides a simple UI for managing consul resources in Consul with more powerful features.

<strong>Note:</strong> This project is currently adding more features. Breaking changes and inavailabilities may occur without notice. Any suggestions and feedback are greatly welcomed.

## Features

Key-Value:

- [x] Tree-like view
- [x] Create new key or folder
- [x] View key value
- [x] Edit key value
- [x] Delete key or folder
- [x] Delete preview
- [ ] Import / Export
- [ ] Jinja2 template support (it's helpful for migration) (maybe in premium version)
- [ ] Config highlight
- [ ] Encrypted import / export
- [x] Modification recording / multiple version control and rollback

ACL Token:

- [x] List tokens
- [x] Create new token
- [x] Create token with exclusive policy 
- [x] View token detail
- [ ] Edit token policies and roles (WIP)
- [x] Delete token
- [x] Delete preview
- [x] Token application

ACL Policy:

- [x] List policies
- [x] Create new policy
- [x] View policy detail
- [ ] Edit policy rule (WIP)
- [x] Delete policy
- [x] Delete preview
- [x] Policy rule CRUD and preview (except service intention)

ACL Role:
- [ ] Role management (WIP)

Admin:

- [x] List notification
- [x] Update notification list automatically
- [x] Review token application request

Other:
- [ ] I18N support (WIP)
- [x] Mobile Equipment UI Support

Deployment:
- [x] Binary File
- [x] Windows Installer
- [ ] Docker Compose (WIP)
- [ ] Helm Chart

## Try Consee

### Using Docker Compose

```shell
git clone https://github.com/FlyingOnion/consee.git
cd consee
docker compose up -d

# consee frontend is available in
# http://localhost:5173
```

Notify that backend is NOT serving static files in dev mode.

For Chinese friends, you can uncomment changing mirrors commands in `frontend/Dockerfile`, `backend/Dockerfile` to speed up building docker mirrors.

### Using Helm Chart (TODO)

## Local Build

You can also build consee from source code. The following dependencies should be installed before building:

- go 1.25+ (`waitgroup.Go` is used)
- bun
- task

Build binary:

```shell
git clone https://github.com/FlyingOnion/consee.git
cd consee
task build-all
```

Prepare config file:
```shell
cd build/bin
mkdir config
cat <<EOF > config/config.yaml
consul:
  address: localhost:8500
  datacenter: dc1
  admin_token: <PASSWORD>
log_level: info
EOF
```

## Config Highlight (WIP)

https://prismjs.com/download.html#themes=prism&languages=markup+clike+javascript+cmake+hcl+ini+json+json5+jsonp+lua+makefile+plsql+properties+qml+sql+toml+yaml&plugins=line-numbers

- cmake
- hcl
- ini
- json, json5, jsonp
- lua
- makefile
- properties
- qml
- sql, plsql
- toml
- xml
- yaml/yml

## Version History

- v0.25.11: Fixed some bugs and improved UI
- v0.25.10: Initial release