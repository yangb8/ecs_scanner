# ecs_scanner


Build
=====

go get github.com/mitchellh/gox

```
cd <target>
gox -output="../bin/{{.Dir}}-{{.OS}}-{{.Arch}}"
cd .. && ls bin
```
