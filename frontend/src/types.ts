export interface Photo {
  id: number;
  title?: string;
  photo_url: string;
  comment?: string;
  metadata?: Metadata;
  tags: Tag[];
}

export interface Metadata {
  altitude: string;
  aperture: string;
  cameramake: string;
  cameramodel: string;
  createdat: string;
  focallength: string;
  imageheight: string;
  imagetype: string;
  imagewidth: string;
  iso: string;
  latitude: string;
  lensmake: string;
  lensmodel: string;
  longitude: string;
  shutterspeed: string;
}

export interface Tag {
  title: string;
  comment: string;
}

export interface PostHead {
  slug: string;
  published: boolean;
  title: string;
  created_at: string;
  comment: string;
  tags: Tag[];
}

export interface Post extends PostHead {
  content: string;
}

export interface ErrorResponse {
  error: string;
}
