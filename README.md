# Scratch부터 만드는 컨테이너

- Liz Rice의 Container from scratch
- 간단한 컨테이너 런타임 만들기
- **cgroup: 그룹 내에서 얼마의 자원을 사용할지에 대한 제어**
- **namespace: 그룹에서 어떤 것들을 볼 수 있는지에 대한 제어**

## 실행 흐름

### RUN

```bash
go run main.go run /bin/bash
```

명령어를 실행하면, `run()`함수가 먼저 호출되어서, 

arg 2번 인덱스 이후에 실제로 전달하는 커맨드(/bin/bash)부분과, 그의 pid가 출력된다.

```bash
Running [/bin/bash] as 62832
```

그 뒤, `cmd`에 실행할 명령어를 정의한다. `/proc/self/exe` 는 현재 프로세스의 실행 바이너리에 대한 심볼릭 링크를 가진다. 

```bash
[root@chaewoon]# ls -l /proc/62832/exe
lrwxrwxrwx. 1 root root 0 Aug  5 13:45 /proc/62832/exe -> /root/.cache/go-build/03/03297bd84654a0183f92872e5322eac5ca84c38280f7ad047e5a9183ffa67b85-d/main
```

그 뒤, 서브프로세스의 표준입출력을 이번 프로세스, 즉 현재 터미널의 표준입출력에 연결한다.

그 뒤, `SysProcAttr`에서 다음의 플래그들을 켜면 된다:

- `NEWUTS`: UTS: Unix Time-Sharing Namespaces
    - 격리된 호스트네임 및 도메인 이름을 가지게 한다
- `NEWPID`: 프로세스를 격리한다(자식 프로세스는 이제 PID 1로 시작한다)
- `NEWNS`: 새로운 마운트 네임스페이스 그룹 분리 및 독립

- `NEWNET`: 새로운 네트워크 네임스페이스 그룹 생성(여기서는 구현되지 않음)
- `NEWUSER`: 새로운 유저 네임스페이스 그룹 생성(여기서는 구현되지 않음)

이제, 자식 프로세스 명령어를 실행한다.

### Child

cgroup설정은 `cg()` 에서 호출된다:

- `sys/fs/cgroup/chaewoon` 디렉토리가 생성
- 이 컨트롤 그룹은 최대 20개의 프로세스만 허용
- 현재 프로세스를 해당 cgroup에 등록

namespace관련 다음의 작업이 일어난다:

- hostname 수정
- 루트 수정
- 루트로 디렉토리 이동
- proc 마운트
    - PID, 메모리 등의 프로세스 정보가 확인 가느아다.

그 뒤, cmd를 실행하는데, 격리되었으므로 내부에서는 PID가 1로 보인다.
```bash
Running [/bin/bash] as 1
```

종료 후, `/proc` 마운트를 해제한다.

## 프로세스 계층 구조

```bash
|\_ bash
|    \_ go run main.go run /bin/bash
|        \_ /root/.cache/go-build/dc/dc4801fa9ce13e84261c2c20041646c3e32e357e77738f7bca94
|            \_ /proc/self/exe child /bin/bash
|                \_ /bin/bash
```

```bash
root@container:/# ps aux
USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root           1  0.0  0.0 1225472 1792 ?        Sl   11:36   0:00 /proc/self/exe child /bin/bash
root           6  0.0  0.0   4320  3456 ?        S    11:36   0:00 /bin/bash
root          10  0.0  0.0   5156  3072 ?        R+   14:04   0:00 ps aux
root@container:/#
```
