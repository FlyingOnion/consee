import mitt from "mitt";
import type { AuthResult } from "./alova";

const emitter = mitt<Event>();

export type Event = {
  login: Pick<AuthResult, 'valid' | 'admin'>;
  logout: void;
  kvPathOpen: string;
  kvPathClose: string;
  kvCreate: void;
  kvDelete: void;
  tokenCreate: void;
  tokenDelete: void;
  policyCreate: void;
  policyDelete: void;
  roleCreate: void;
  roleDelete: void;
  openNotificationsChange: number;
  notificationResolve: void;
};

export default emitter;
