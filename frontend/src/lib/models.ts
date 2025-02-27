export interface File {
  fileID: number;
  name: string;
  path: string;
  owner: number;
  parentID: number;

  header: string;
  key: string;
  nonce: string;
}

export interface Folder {
  folderID: number;
  name: string;
  path: string;
  owner: number;
  parentID: number;
}

export interface FolderContent {
  files: File[];
  folders: Folder[];
}

export interface User {
  userID: number;
  username: string;
  masterKeySalt: string;

  shareKeyPublic: string;
  shareKeyPrivate: string;
  shareKeyNonce: string;
}
