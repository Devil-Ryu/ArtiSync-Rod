// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {context} from '../models';
import {utils} from '../models';
import {proto} from '../models';

export function CheckAuthentication():Promise<{[key: string]: string}>;

export function CheckConfig(arg1:any):Promise<void>;

export function Init(arg1:context.Context,arg2:utils.Article,arg3:number):Promise<void>;

export function LoadConfig(arg1:{[key: string]: any},arg2:boolean):Promise<{[key: string]: any}>;

export function LoadCookies():Promise<Array<proto.NetworkCookieParam>>;

export function Login():Promise<void>;

export function OpenPage(arg1:string):Promise<void>;

export function RUN():Promise<void>;

export function Run():Promise<void>;

export function SetConfig(arg1:boolean):Promise<void>;

export function SetCookies(arg1:Array<proto.NetworkCookieParam>):Promise<void>;

export function UpdatePlatformInfo():Promise<void>;
