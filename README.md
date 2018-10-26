# run

run은 원하는 환경변수를 설정 한 뒤 특정 프로그램을 실행시키는 프로그램 입니다.
OS에서 설정된 환경변수를 불러온 후 수정하거나 다른 환경변수를 추가한 후
프로그램을 실행시키게 됩니다.


### 실행

run을 통해 다음처럼 HOUDINI_PATH 환경변수를 직접 설정한 후 houdini를 실행할 수 있습니다.

```
run -env="HOUDINI_PATH=/vfx/global/houdini" houdini
```

또는 후디니 관련 환경변수 설정이 들어가 있는 env 파일을 통해서도 환경변수를 설정할 수 있습니다.

```
run -envfile="/vfx/global/houdini/env/all.env"
```

두 플래그를 모두 사용하면 -env 플래그의 값이 먼저 설정됩니다.

즉, 아래 예제에서 HOUDINI_PATH가 먼저 설정되고 그다음 all.env안의 환경변수들이 설정됩니다.

```
run -env="HOUDINI_PATH=/vfx/global/houdini" -envfile="/vfx/global/houdini/env/all.env"
```

두 플래그 모두 콤마(,)를 사용하여 여러 환경변수(파일)을 추가할 수 있습니다.

```
run -env="HOUDINI_PATH=/vfx/global/houdini,RENDERMAN_PATH=/vfx/global/renderman" -envfile="/vfx/global/houdini/env/all.env,/vfx/global/houdini/env/myteam.env"
```

-envfile 플래그 경로 뒤에 물음표가 붙으면 그 env이 없어도 에러가 나지 않습니다.
즉 아래에서 maybe.env파일이 없다고 에러가 나지 않습니다.

```
run -envfile="must.env,maybe.env?"
```


### env 파일

env 파일 형식은 기본적으로 bash 스타일의 환경변수 설정을 따르되
export 키워드는 사용하지 않습니다.

```
HOUDINI_PATH=/vfx/global/houdini
NUKE_PATH=/vfx/global/nuke
MAYA_PATH=/vfx/global/maya
```

#### 환경변수 치환

하나의 환경변수를 설정할 때 다른 환경변수를 사용할 수 있으며
이 때는 $를 환경변수 이름 앞에 붙여 사용하면 됩니다.

```
GLOBAL_PATH=/vfx/global
HOUDINI_PATH=$GLOBAL_PATH/houdini
NUKE_PATH=$GLOBAL_PATH/nuke
MAYA_PATH=$GLOBAL_PATH/maya
```

이는 OS에서 미리 설정된 환경변수도 적용됩니다.

```
# GLOBAL_PATH를 OS에서 설정한 후,
HOUDINI_PATH=$GLOBAL_PATH/houdini
NUKE_PATH=$GLOBAL_PATH/nuke
MAYA_PATH=$GLOBAL_PATH/maya
```

#### OS별 문자열 자동 변경

환경변수 경로명과 구분자는은 리눅스 형식으로 작성하고 run 실행 환경이 윈도우즈 일 때는
자동으로 윈도우즈 스타일로 변경됩니다.

윈도우즈에서 변경되는 문자는 현재 2가지 입니다.

```
경로 구분자  :	/ -> \
리스트 구분자:	: -> ;
```

다만 백쿼트(\`)로 감싸진 문자열은 변경되지 않습니다. 즉,

```
HOUDINI_PATH=`Z:`/global/houdini
HOUDINI_PATH=`C:`/my/houdini:$HOUDINI_PATH
```

처럼 설정한 후 윈도우즈에서 run을 실행시키면 HOUDINI_PATH는
`C:\my\houdini;Z:\global\houdini` 가 됩니다.

#### 주석

`#` 으로 시작하는 줄은 주석입니다.

```
# test.env
#
# 환경변수 테스트
TEST=a

# 또다른 환경변수
ITS_LIKE=tt
```


