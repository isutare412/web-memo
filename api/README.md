# APIs

| Name                   | Method | URI                                                     |   Authorization    |
| ---------------------- | ------ | ------------------------------------------------------- | :----------------: |
| googleSignIn           | GET    | /api/v1/google/sign-in                                  |         N          |
| googleSignInFinish     | GET    | /api/v1/google/sign-in/finish                           |         N          |
| getSelfUser            | GET    | /api/v1/users/me                                        |         Y          |
| refreshUserToken       | POST   | /api/v1/users/refresh-token
| signOutUser            | GET    | /api/v1/users/sign-out                                  |         N          |
| getMemo                | GET    | /api/v1/memos/{memoID}                                  | Y (N if published) |
| listMemos              | GET    | /api/v1/memos                                           |         Y          |
| createMemo             | POST   | /api/v1/memos                                           |         Y          |
| replaceMemo            | PUT    | /api/v1/memos/{memoID}                                  |         Y          |
| publishMemo            | POST   | /api/v1/memos/{memoID}/publish                          |         Y          |
| deleteMemo             | DELETE | /api/v1/memos/{memoID}                                  |         Y          |
| listMemoTags           | GET    | /api/v1/memos/{memoID}/tags                             | Y (N if published) |
| replaceMemoTags        | PUT    | /api/v1/memos/{memoID}/tags                             |         Y          |
| getCollaborator        | GET    | /api/v1/memos/{memoID}/collaborators/{userID}           |         Y          |
| listCollaborators      | GET    | /api/v1/memos/{memoID}/collaborators                    |         Y          |
| requestCollaboration   | PUT    | /api/v1/memos/{memoID}/collaborators/{userID}           |         Y          |
| authorizeCollaboration | POST   | /api/v1/memos/{memoID}/collaborators/{userID}/authorize |         Y          |
| cancelCollaboration    | DELETE | /api/v1/memos/{memoID}/collaborators/{userID}           |         Y          |
| getSubscriber          | GET    | /api/v1/memos/{memoID}/subscribers/{userID}             |         Y          |
| listSubscribers        | GET    | /api/v1/memos/{memoID}/subscribers                      |         Y          |
| subscribeMemo          | PUT    | /api/v1/memos/{memoID}/subscribers/{userID}             |         Y          |
| unsubscribeMemo        | DELETE | /api/v1/memos/{memoID}/subscribers/{userID}             |         Y          |
| listTags               | GET    | /api/v1/tags                                            |         Y          |

```
Authorization: Bearer <idToken>
```

## Using cookie
```
Cookie: wmToken=<idToken>
```
