export interface UserSearchFormModel {
  created_at: string[];
  dept_id?: number;
  email: string;
  phone: string;
  post_id?: number;
  role_id?: number;
  status?: number;
  user_type?: string;
  username: string;
}

export interface UserListItem {
  [key: string]: any;
  dashboard?: string;
  dept_name?: string;
  email?: string;
  id: number;
  nickname?: string;
  phone?: string;
  post_name?: string;
  role_name?: string;
  status?: number;
  user_type?: string;
  username: string;
}

export interface ColumnOptionItem {
  label: string;
  value: string;
}
