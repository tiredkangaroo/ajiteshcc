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
