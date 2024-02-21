# APIs

| Name               | Method | URI                                |   Authorization    |
| ------------------ | ------ | ---------------------------------- | :----------------: |
| googleSignIn       | GET    | /api/v1/google/sign-in             |         N          |
| googleSignInFinish | GET    | /api/v1/google/sign-in/finish      |         N          |
| getSelfUser        | GET    | /api/v1/users/me                   |         Y          |
| signOutUser        | GET    | /api/v1/users/sign-out             |         N          |
| getMemo            | GET    | /api/v1/memos/{memoID}             | Y (N if published) |
| listMemos          | GET    | /api/v1/memos                      |         Y          |
| createMemo         | POST   | /api/v1/memos                      |         Y          |
| replaceMemo        | PUT    | /api/v1/memos/{memoID}             |         Y          |
| publishMemo        | POST   | /api/v1/memos/{memoID}/publish     |         Y          |
| subscribeMemo      | POST   | /api/v1/memos/{memoID}/subscribe   |         Y          |
| deleteMemo         | DELETE | /api/v1/memos/{memoID}             |         Y          |
| listMemoTags       | GET    | /api/v1/memos/{memoID}/tags        | Y (N if published) |
| replaceMemoTags    | PUT    | /api/v1/memos/{memoID}/tags        |         Y          |
| listSubscribers    | GET    | /api/v1/memos/{memoID}/subscribers |         Y          |
| listTags           | GET    | /api/v1/tags                       |         Y          |

```
Authorization: Bearer <idToken>
```

## Using cookie
```
Cookie: wmToken=<idToken>
```
