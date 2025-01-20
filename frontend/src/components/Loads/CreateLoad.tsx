//src/components/Loads/CreateLoad.tsx

import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  TextField,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material';
import Grid2 from '@mui/material/Grid2';
import { ExpandMore as ExpandMoreIcon, ArrowBack } from '@mui/icons-material';
import { loadService } from '../../services/api';
import { Load } from '../../types/load.types';

interface FormSectionProps {
  title: string;
  children: React.ReactNode;
  defaultExpanded?: boolean;
}

type NestedKeyOf<ObjectType extends object> = {
  [Key in keyof ObjectType & (string | number)]: ObjectType[Key] extends object
    ? `${Key}` | `${Key}.${NestedKeyOf<ObjectType[Key]>}`
    : `${Key}`;
}[keyof ObjectType & (string | number)];

const FormSection: React.FC<FormSectionProps> = ({ 
  title, 
  children, 
  defaultExpanded = false 
}) => (
  <Accordion defaultExpanded={defaultExpanded} sx={{ mb: 2 }}>
    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
      <Typography variant="h6">{title}</Typography>
    </AccordionSummary>
    <AccordionDetails>
      <Grid2 container spacing={2}>
        {children}
      </Grid2>
    </AccordionDetails>
  </Accordion>
);

const initialFormData: Partial<Load> = {
  externalTMSLoadID: "",
  freightLoadID: "",
  status: {
    code: {
      key: "2102",
      value: "Covered"
    },
    notes: "",
    description: ""
  },
  customer: {
    name: "",
    accountNumber: "",
    contact: {
      name: "",
      email: "",
      phone: ""
    },
    address: {
      street: "",
      city: "",
      state: "",
      zipCode: ""
    }
  },
  billTo: {
    name: "",
    accountNumber: "",
    contact: {
      name: "",
      email: "",
      phone: ""
    },
    address: {
      street: "",
      city: "",
      state: "",
      zipCode: ""
    }
  },
  pickup: {
    facilityName: "",
    scheduledTime: "",
    contact: {
      name: "",
      phone: ""
    },
    address: {
      street: "",
      city: "",
      state: "",
      zipCode: ""
    }
  },
  consignee: {
    facilityName: "",
    scheduledTime: "",
    contact: {
      name: "",
      phone: ""
    },
    address: {
      street: "",
      city: "",
      state: "",
      zipCode: ""
    }
  },
  carrier: {
    name: "",
    scac: "",
    contact: {
      name: "",
      phone: ""
    },
    equipment: {
      type: "DryVan",
      length: "53"
    }
  },
  rateData: {
    baseRate: 0,
    fuelSurcharge: 0,
    totalRate: 0,
    currency: "USD"
  },
  specifications: {
    temperature: {
      min: 35,
      max: 75,
      unit: "F"
    },
    serviceLevel: "Standard",
    specialInstructions: ""
  },
  inPalletCount: 0,
  outPalletCount: 0,
  numCommodities: 0,
  totalWeight: 0,
  billableWeight: 0,
  poNums: "",
  operator: "SYSTEM",
  routeMiles: 0
};

export const CreateLoad = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<Partial<Load>>(initialFormData);

  const handleInputChange = (path: NestedKeyOf<Load>, value: string | number) => {
    setFormData(prev => {
      const keys = path.split('.');
      const newData = { ...prev };
      let current: any = newData;
      
      for (let i = 0; i < keys.length - 1; i++) {
        current = current[keys[i]];
      }
      
      current[keys[keys.length - 1]] = value;
      return newData;
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setLoading(true);
      const response = await loadService.createLoad(formData as Omit<Load, 'id' | 'createdAt' | 'updatedAt'>);
      navigate(`/loads/${response.id}`);
    } catch (error) {
      console.error('Error creating load:', error);
      // TODO: Add error handling
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit}>
      <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', gap: 2 }}>
        <Button
          startIcon={<ArrowBack />}
          onClick={() => navigate('/')}
          variant="outlined"
        >
          Back to Loads
        </Button>
        <Typography variant="h5">Create New Load</Typography>
      </Box>

      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        {/* Basic Information */}
        <FormSection title="Basic Information" defaultExpanded>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="External TMS Load ID"
              value={formData.externalTMSLoadID}
              onChange={(e) => handleInputChange('externalTMSLoadID', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Freight Load ID"
              value={formData.freightLoadID}
              onChange={(e) => handleInputChange('freightLoadID', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Status Code Key"
              value={formData.status?.code.key}
              onChange={(e) => handleInputChange('status.code.key', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Status Code Value"
              value={formData.status?.code.value}
              onChange={(e) => handleInputChange('status.code.value', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              label="PO Numbers"
              value={formData.poNums}
              onChange={(e) => handleInputChange('poNums', e.target.value)}
              helperText="Separate multiple PO numbers with commas"
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              label="Operator"
              value={formData.operator}
              onChange={(e) => handleInputChange('operator', e.target.value)}
            />
          </Grid2>
        </FormSection>

        {/* Customer Information */}
        <FormSection title="Customer Information">
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Customer Name"
              value={formData.customer?.name}
              onChange={(e) => handleInputChange('customer.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Account Number"
              value={formData.customer?.accountNumber}
              onChange={(e) => handleInputChange('customer.accountNumber', e.target.value)}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              label="Contact Name"
              value={formData.customer?.contact.name}
              onChange={(e) => handleInputChange('customer.contact.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              label="Contact Email"
              value={formData.customer?.contact.email}
              onChange={(e) => handleInputChange('customer.contact.email', e.target.value)}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              label="Contact Phone"
              value={formData.customer?.contact.phone}
              onChange={(e) => handleInputChange('customer.contact.phone', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Street Address"
              value={formData.customer?.address.street}
              onChange={(e) => handleInputChange('customer.address.street', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="City"
              value={formData.customer?.address.city}
              onChange={(e) => handleInputChange('customer.address.city', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="State"
              value={formData.customer?.address.state}
              onChange={(e) => handleInputChange('customer.address.state', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Zip Code"
              value={formData.customer?.address.zipCode}
              onChange={(e) => handleInputChange('customer.address.zipCode', e.target.value)}
            />
          </Grid2>
        </FormSection>

        {/* Bill To Information */}
        <FormSection title="Bill To Information">
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Bill To Name"
              value={formData.billTo?.name}
              onChange={(e) => handleInputChange('billTo.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Bill To Account Number"
              value={formData.billTo?.accountNumber}
              onChange={(e) => handleInputChange('billTo.accountNumber', e.target.value)}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              label="Contact Name"
              value={formData.billTo?.contact.name}
              onChange={(e) => handleInputChange('billTo.contact.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              label="Contact Email"
              value={formData.billTo?.contact.email}
              onChange={(e) => handleInputChange('billTo.contact.email', e.target.value)}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              label="Contact Phone"
              value={formData.billTo?.contact.phone}
              onChange={(e) => handleInputChange('billTo.contact.phone', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Street Address"
              value={formData.billTo?.address.street}
              onChange={(e) => handleInputChange('billTo.address.street', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="City"
              value={formData.billTo?.address.city}
              onChange={(e) => handleInputChange('billTo.address.city', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="State"
              value={formData.billTo?.address.state}
              onChange={(e) => handleInputChange('billTo.address.state', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Zip Code"
              value={formData.billTo?.address.zipCode}
              onChange={(e) => handleInputChange('billTo.address.zipCode', e.target.value)}
            />
          </Grid2>
        </FormSection>

        {/* Pickup Information */}
        <FormSection title="Pickup Information">
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Facility Name"
              value={formData.pickup?.facilityName}
              onChange={(e) => handleInputChange('pickup.facilityName', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              type="datetime-local"
              label="Scheduled Time"
              value={formData.pickup?.scheduledTime?.split('Z')[0]}
              onChange={(e) => handleInputChange('pickup.scheduledTime', `${e.target.value}Z`)}
              InputLabelProps={{ shrink: true }}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Contact Name"
              value={formData.pickup?.contact.name}
              onChange={(e) => handleInputChange('pickup.contact.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Contact Phone"
              value={formData.pickup?.contact.phone}
              onChange={(e) => handleInputChange('pickup.contact.phone', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Street Address"
              value={formData.pickup?.address.street}
              onChange={(e) => handleInputChange('pickup.address.street', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="City"
              value={formData.pickup?.address.city}
              onChange={(e) => handleInputChange('pickup.address.city', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="State"
              value={formData.pickup?.address.state}
              onChange={(e) => handleInputChange('pickup.address.state', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Zip Code"
              value={formData.pickup?.address.zipCode}
              onChange={(e) => handleInputChange('pickup.address.zipCode', e.target.value)}
            />
          </Grid2>
        </FormSection>
        <FormSection title="Delivery Information">
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Facility Name"
              value={formData.consignee?.facilityName}
              onChange={(e) => handleInputChange('consignee.facilityName', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              type="datetime-local"
              label="Scheduled Time"
              value={formData.consignee?.scheduledTime?.split('Z')[0]}
              onChange={(e) => handleInputChange('consignee.scheduledTime', `${e.target.value}Z`)}
              InputLabelProps={{ shrink: true }}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Contact Name"
              value={formData.consignee?.contact.name}
              onChange={(e) => handleInputChange('consignee.contact.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Contact Phone"
              value={formData.consignee?.contact.phone}
              onChange={(e) => handleInputChange('consignee.contact.phone', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Street Address"
              value={formData.consignee?.address.street}
              onChange={(e) => handleInputChange('consignee.address.street', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="City"
              value={formData.consignee?.address.city}
              onChange={(e) => handleInputChange('consignee.address.city', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="State"
              value={formData.consignee?.address.state}
              onChange={(e) => handleInputChange('consignee.address.state', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Zip Code"
              value={formData.consignee?.address.zipCode}
              onChange={(e) => handleInputChange('consignee.address.zipCode', e.target.value)}
            />
          </Grid2>
        </FormSection>

        {/* Carrier Information */}
        <FormSection title="Carrier Information">
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Carrier Name"
              value={formData.carrier?.name}
              onChange={(e) => handleInputChange('carrier.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="SCAC"
              value={formData.carrier?.scac}
              onChange={(e) => handleInputChange('carrier.scac', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Contact Name"
              value={formData.carrier?.contact.name}
              onChange={(e) => handleInputChange('carrier.contact.name', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Contact Phone"
              value={formData.carrier?.contact.phone}
              onChange={(e) => handleInputChange('carrier.contact.phone', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Equipment Type"
              value={formData.carrier?.equipment.type}
              onChange={(e) => handleInputChange('carrier.equipment.type', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Equipment Length"
              value={formData.carrier?.equipment.length}
              onChange={(e) => handleInputChange('carrier.equipment.length', e.target.value)}
            />
          </Grid2>
        </FormSection>

        {/* Rate Information */}
        <FormSection title="Rate Information">
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Base Rate"
              value={formData.rateData?.baseRate}
              onChange={(e) => handleInputChange('rateData.baseRate', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Fuel Surcharge"
              value={formData.rateData?.fuelSurcharge}
              onChange={(e) => handleInputChange('rateData.fuelSurcharge', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Total Rate"
              value={formData.rateData?.totalRate}
              onChange={(e) => handleInputChange('rateData.totalRate', Number(e.target.value))}
            />
          </Grid2>
        </FormSection>

        {/* Specifications */}
        <FormSection title="Specifications">
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Temperature Min"
              value={formData.specifications?.temperature.min}
              onChange={(e) => handleInputChange('specifications.temperature.min', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Temperature Max"
              value={formData.specifications?.temperature.max}
              onChange={(e) => handleInputChange('specifications.temperature.max', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              label="Temperature Unit"
              value={formData.specifications?.temperature.unit}
              onChange={(e) => handleInputChange('specifications.temperature.unit', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              required
              label="Service Level"
              value={formData.specifications?.serviceLevel}
              onChange={(e) => handleInputChange('specifications.serviceLevel', e.target.value)}
            />
          </Grid2>
          <Grid2 size={6}>
            <TextField
              fullWidth
              multiline
              rows={2}
              label="Special Instructions"
              value={formData.specifications?.specialInstructions}
              onChange={(e) => handleInputChange('specifications.specialInstructions', e.target.value)}
            />
          </Grid2>
        </FormSection>

        {/* Quantities */}
        <FormSection title="Quantities">
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="In Pallet Count"
              value={formData.inPalletCount}
              onChange={(e) => handleInputChange('inPalletCount', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Out Pallet Count"
              value={formData.outPalletCount}
              onChange={(e) => handleInputChange('outPalletCount', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Number of Commodities"
              value={formData.numCommodities}
              onChange={(e) => handleInputChange('numCommodities', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Total Weight"
              value={formData.totalWeight}
              onChange={(e) => handleInputChange('totalWeight', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Billable Weight"
              value={formData.billableWeight}
              onChange={(e) => handleInputChange('billableWeight', Number(e.target.value))}
            />
          </Grid2>
          <Grid2 size={4}>
            <TextField
              fullWidth
              required
              type="number"
              label="Route Miles"
              value={formData.routeMiles}
              onChange={(e) => handleInputChange('routeMiles', Number(e.target.value))}
            />
          </Grid2>
        </FormSection>

        <Box sx={{ mt: 2, display: 'flex', justifyContent: 'flex-end', gap: 2 }}>
          <Button
            variant="outlined"
            onClick={() => navigate('/')}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            type="submit"
            disabled={loading}
          >
            Create Load
          </Button>
        </Box>
      </Box>
    </Box>
  );
};