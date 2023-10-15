# APIs

| Name               | Method | URI                           | Authorization |
| ------------------ | ------ | ----------------------------- | :-----------: |
| googleSignIn       | GET    | /api/v1/google/sign-in        |       N       |
| googleSignInFinish | GET    | /api/v1/google/sign-in/finish |       N       |
| getSelfUser        | GET    | /api/v1/users/me              |       Y       |
| signOutUser        | GET    | /api/v1/users/sign-out        |       Y       |
| getMemo            | GET    | /api/v1/memos/{memoID}        |       Y       |
| listMemos          | GET    | /api/v1/memos                 |       Y       |
| createMemo         | POST   | /api/v1/memos                 |       Y       |
| replaceMemo        | PUT    | /api/v1/memos/{memoID}        |       Y       |
| deleteMemo         | DELETE | /api/v1/memos/{memoID}        |       Y       |
| listTags           | GET    | /api/v1/memos/{memoID}/tags   |       Y       |
| replaceTags        | PUT    | /api/v1/memos/{memoID}/tags   |       Y       |

# Authorization

## Using header

```
Authorization: Bearer <idToken>
```

## Using cookie
```
Cookie: wmToken=<idToken>
```
