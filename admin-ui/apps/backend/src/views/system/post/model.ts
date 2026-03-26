export interface PostSearchFormModel {
  code: string;
  created_at: string[];
  name: string;
  status?: number;
}

export interface PostListItem {
  [key: string]: any;
  code: string;
  id: number;
  name: string;
  remark?: string;
  sort?: number;
  status?: number;
}

export interface PostFormModel {
  code: string;
  id?: number;
  name: string;
  remark: string;
  sort: number;
  status: number;
}
