// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {controller} from '../models';
import {db} from '../models';
import {utils} from '../models';

export function CheckAuthentication():Promise<{[key: string]: string}>;

export function CheckConfig(arg1:any):Promise<void>;

export function HasAccount():Promise<boolean>;

export function HasDBController():Promise<boolean>;

export function HasRODController():Promise<boolean>;

export function InitRod(arg1:controller.DBController,arg2:controller.RODController,arg3:any):Promise<void>;

export function LoadConfig():Promise<{[key: string]: any}>;

export function Login():Promise<Array<db.Account>>;

export function OpenPage(arg1:string):Promise<void>;

export function Publish():Promise<void>;

export function SetAccount(arg1:db.Account):Promise<void>;

export function SetArticle(arg1:utils.Article):Promise<void>;

export function SetController(arg1:controller.DBController,arg2:controller.RODController):Promise<void>;

export function Start(arg1:controller.DBController,arg2:controller.RODController,arg3:any,arg4:db.Account,arg5:utils.Article,arg6:any):Promise<void>;

export function UpdatePlatformInfo():Promise<void>;
