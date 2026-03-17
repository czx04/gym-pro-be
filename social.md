# Social Mobile Integration Guide

Tai lieu nay danh cho mobile team (React Native/Expo) de goi dung flow Social da implement tren backend.

## 1) Trang thai hien tai

Flow upload media da full end-to-end:

1. Xin signature tu backend.
2. Upload truc tiep len Cloudinary.
3. Confirm media voi backend.
4. Tao post va attach media vao post theo transaction.
5. Doc feed/post detail/user posts co media.

Tat ca endpoint social deu can Bearer token.

## 2) Endpoint can dung cho mobile

- GET /api/v1/social/feed
- POST /api/v1/social/posts
- GET /api/v1/social/posts/:postId
- GET /api/v1/social/users/:userId/profile
- GET /api/v1/social/users/:userId/posts
- POST /api/v1/social/media/signature
- POST /api/v1/social/media/confirm

## 3) Flow tao bai viet co anh/video

### Buoc 1: Xin signature

POST /api/v1/social/media/signature

Request:

{
"resource_type": "image",
"folder": "posts/{userId}"
}

Response data quan trong:

- cloud_name
- api_key
- timestamp
- folder
- public_id
- signature
- upload_url
- expires_in
- resource_type

Ghi chu:

- folder phai bat dau bang posts/{authenticatedUserId}
- resource_type: image hoac video
- expires_in hien tai la 300s

### Buoc 2: Upload thang len Cloudinary

Client KHONG upload file qua backend.

Field multipart can gui len upload_url:

- file
- api_key
- timestamp
- signature
- folder
- public_id

Cloudinary se tra ve it nhat:

- public_id
- secure_url
- resource_type
- bytes

### Buoc 3: Confirm media voi backend

POST /api/v1/social/media/confirm

Request:

{
"public_id": "posts/{userId}/{generated_id}",
"secure_url": "https://res.cloudinary.com/...",
"resource_type": "image",
"bytes": 482193
}

Expected:

- asset_state = ready

Neu public_id khong thuoc user dang login, backend tra forbidden.

### Buoc 4: Tao post va attach media

POST /api/v1/social/posts

Request toi thieu:

{
"caption": "Hom nay leg day",
"media": [
{
"public_id": "posts/{userId}/{generated_id}",
"resource_type": "image"
}
]
}

Logic server:

- Tao post.
- Kiem tra tung media asset phai la ready + dung owner.
- Attach media vao post_media theo order_index.
- Cap nhat asset status attached.
- Tat ca trong 1 transaction, loi thi rollback.

## 4) Doc du lieu de render UI

### GET /api/v1/social/feed

Dung cho Community list.

### GET /api/v1/social/posts/:postId

Dung cho ViewPostModal.

### GET /api/v1/social/users/:userId/posts

Dung cho UserProfileModal post list.

Post output da co media array:

{
"id": "...",
"caption": "...",
"media": [
{
"type": "image",
"url": "https://res.cloudinary.com/..."
}
]
}

## 5) Mapping de mobile implement nhanh

Sequence trong CreatePostModal:

1. call /social/media/signature
2. upload Cloudinary
3. call /social/media/confirm
4. call /social/posts
5. invalidate query feed + user posts

## 6) Error handling khuyen nghi

Case thuong gap:

- 403 invalid folder/public_id ownership
- 404 media asset not found (confirm sai public_id)
- 409 media asset not ready or not owned khi create post
- 422 request body/field khong hop le

Mobile nen:

- show message than thien
- cho retry o buoc upload/confirm
- neu create post fail 409, force re-confirm media hoac upload lai

## 7) Pseudo-code React Native

```ts
async function createPostWithMedia(
  file: FileLike,
  caption: string,
  userId: string,
) {
  const sig = await api.post("/social/media/signature", {
    resource_type: "image",
    folder: `posts/${userId}`,
  });

  const form = new FormData();
  form.append("file", file as any);
  form.append("api_key", sig.data.api_key);
  form.append("timestamp", String(sig.data.timestamp));
  form.append("signature", sig.data.signature);
  form.append("folder", sig.data.folder);
  form.append("public_id", sig.data.public_id);

  const cloudRes = await fetch(sig.data.upload_url, {
    method: "POST",
    body: form,
  }).then((r) => r.json());

  await api.post("/social/media/confirm", {
    public_id: cloudRes.public_id,
    secure_url: cloudRes.secure_url,
    resource_type: cloudRes.resource_type,
    bytes: cloudRes.bytes,
  });

  const post = await api.post("/social/posts", {
    caption,
    media: [
      {
        public_id: cloudRes.public_id,
        resource_type: cloudRes.resource_type,
      },
    ],
  });

  return post.data;
}
```
