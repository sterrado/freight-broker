// src/types/load.types.ts

export interface Address {
    city: string;
    state: string;
    street: string;
    zipCode: string;
  }
  
  export interface Contact {
    email?: string;
    name: string;
    phone: string;
  }
  
  export interface StatusCode {
    key: string;
    value: string;
  }
  
  export interface Status {
    code: StatusCode;
    notes: string;
    description: string;
  }
  
  export interface Customer {
    accountNumber: string;
    address: Address;
    contact: Contact;
    name: string;
  }
  
  export interface Location {
    address: Address;
    contact: Contact;
    facilityName: string;
    scheduledTime: string;
  }
  
  export interface Equipment {
    length: string;
    type: string;
  }
  
  export interface Carrier {
    contact: Contact;
    equipment: Equipment;
    name: string;
    scac: string;
  }
  
  export interface RateData {
    baseRate: number;
    currency: string;
    fuelSurcharge: number;
    totalRate: number;
  }
  
  export interface Temperature {
    max: number;
    min: number;
    unit: string;
  }
  
  export interface Specifications {
    serviceLevel: string;
    specialInstructions: string;
    temperature: Temperature;
  }
  
  export interface Load {
    id: string;
    externalTMSLoadID: string;
    freightLoadID: string;
    status: Status;
    customer: Customer;
    billTo: Customer;
    pickup: Location;
    consignee: Location;
    carrier: Carrier;
    rateData: RateData;
    specifications: Specifications;
    inPalletCount: number;
    outPalletCount: number;
    numCommodities: number;
    totalWeight: number;
    billableWeight: number;
    poNums: string;
    operator: string;
    routeMiles: number;
    createdAt: string;
    updatedAt: string;
  }
  
  export interface LoadsResponse {
    loads: Load[];
    total: number;
    page: number;
    size: number;
  }