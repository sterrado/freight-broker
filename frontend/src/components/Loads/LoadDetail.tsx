// src/components/Loads/LoadDetail.tsx

import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Paper,
  Typography,
  Button,
  Grid2,
  CircularProgress,
} from '@mui/material';
import { ArrowBack } from '@mui/icons-material';
import { loadService } from '../../services/api';
import { Load } from '../../types/load.types';

interface DetailSectionProps {
  title: string;
  children: React.ReactNode;
}

const DetailSection: React.FC<DetailSectionProps> = ({ title, children }) => (
  <Box sx={{ mb: 4 }}>
    <Typography variant="h6" sx={{ mb: 2 }}>
      {title}
    </Typography>
    <Paper sx={{ p: 2 }}>{children}</Paper>
  </Box>
);

interface InfoRowProps {
  label: string;
  value: string | number | undefined;
}

const InfoRow: React.FC<InfoRowProps> = ({ label, value }) => (
    <Grid2 container spacing={2}>
      <Grid2 size={4}>
        <Typography color="textSecondary">{label}:</Typography>
      </Grid2>
      <Grid2 size={8}>
        <Typography>{value || 'N/A'}</Typography>
      </Grid2>
    </Grid2>
  );

export const LoadDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [load, setLoad] = useState<Load | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLoad = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const loadData = await loadService.getLoadById(id);
        setLoad(loadData);
      } catch (error) {
        console.error('Error fetching load:', error);
        // TODO: Add error handling
      } finally {
        setLoading(false);
      }
    };

    fetchLoad();
  }, [id]);

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!load) {
    return <Typography>Load not found</Typography>;
  }

  return (
    <Box>
      <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', gap: 2 }}>
        <Button
          startIcon={<ArrowBack />}
          onClick={() => navigate('/')}
          variant="outlined"
        >
          Back to Loads
        </Button>
        <Typography variant="h5">Load Details</Typography>
      </Box>

      <DetailSection title="Basic Information">
        <Grid2 container spacing={2}>
          <Grid2 size={6}>
            <InfoRow label="Freight Load ID" value={load.freightLoadID} />
            <InfoRow label="External TMS ID" value={load.externalTMSLoadID} />
            <InfoRow label="Status" value={load.status.code.value} />
            <InfoRow label="Created At" value={new Date(load.createdAt).toLocaleString()} />
          </Grid2>
          <Grid2 size={6}>
            <InfoRow label="Route Miles" value={`${load.routeMiles} miles`} />
            <InfoRow label="Total Weight" value={`${load.totalWeight} lbs`} />
            <InfoRow label="Billable Weight" value={`${load.billableWeight} lbs`} />
            <InfoRow label="PO Numbers" value={load.poNums} />
          </Grid2>
        </Grid2>
      </DetailSection>

      <DetailSection title="Customer Information">
        <Grid2 container spacing={2}>
          <Grid2 size={6}>
            <InfoRow label="Name" value={load.customer.name} />
            <InfoRow label="Account Number" value={load.customer.accountNumber} />
            <InfoRow label="Contact Name" value={load.customer.contact.name} />
            <InfoRow label="Contact Phone" value={load.customer.contact.phone} />
          </Grid2>
          <Grid2 size={6}>
            <InfoRow label="Address" value={load.customer.address.street} />
            <InfoRow label="City" value={load.customer.address.city} />
            <InfoRow label="State" value={load.customer.address.state} />
            <InfoRow label="Zip Code" value={load.customer.address.zipCode} />
          </Grid2>
        </Grid2>
      </DetailSection>

      <DetailSection title="Pickup Information">
        <Grid2 container spacing={2}>
          <Grid2 size={6}>
            <InfoRow label="Facility Name" value={load.pickup.facilityName} />
            <InfoRow label="Contact Name" value={load.pickup.contact.name} />
            <InfoRow label="Contact Phone" value={load.pickup.contact.phone} />
            <InfoRow 
              label="Scheduled Time" 
              value={new Date(load.pickup.scheduledTime).toLocaleString()} 
            />
          </Grid2>
          <Grid2 size={6}>
            <InfoRow label="Address" value={load.pickup.address.street} />
            <InfoRow label="City" value={load.pickup.address.city} />
            <InfoRow label="State" value={load.pickup.address.state} />
            <InfoRow label="Zip Code" value={load.pickup.address.zipCode} />
          </Grid2>
        </Grid2>
      </DetailSection>

      <DetailSection title="Delivery Information">
        <Grid2 container spacing={2}>
          <Grid2 size={6}>
            <InfoRow label="Facility Name" value={load.consignee.facilityName} />
            <InfoRow label="Contact Name" value={load.consignee.contact.name} />
            <InfoRow label="Contact Phone" value={load.consignee.contact.phone} />
            <InfoRow 
              label="Scheduled Time" 
              value={new Date(load.consignee.scheduledTime).toLocaleString()} 
            />
          </Grid2>
          <Grid2 size={6}>
            <InfoRow label="Address" value={load.consignee.address.street} />
            <InfoRow label="City" value={load.consignee.address.city} />
            <InfoRow label="State" value={load.consignee.address.state} />
            <InfoRow label="Zip Code" value={load.consignee.address.zipCode} />
          </Grid2>
        </Grid2>
      </DetailSection>

      <DetailSection title="Rate Information">
        <Grid2 container spacing={2}>
          <Grid2 size={6}>
            <InfoRow 
              label="Base Rate" 
              value={`$${load.rateData.baseRate.toLocaleString('en-US', { minimumFractionDigits: 2 })}`} 
            />
            <InfoRow 
              label="Fuel Surcharge" 
              value={`$${load.rateData.fuelSurcharge.toLocaleString('en-US', { minimumFractionDigits: 2 })}`} 
            />
            <InfoRow 
              label="Total Rate" 
              value={`$${load.rateData.totalRate.toLocaleString('en-US', { minimumFractionDigits: 2 })}`} 
            />
          </Grid2>
          <Grid2 size={6}>
            <InfoRow label="Currency" value={load.rateData.currency} />
          </Grid2>
        </Grid2>
      </DetailSection>

      <DetailSection title="Carrier Information">
        <Grid2 container spacing={2}>
          <Grid2 size={6}>
            <InfoRow label="Name" value={load.carrier.name} />
            <InfoRow label="SCAC" value={load.carrier.scac} />
            <InfoRow label="Contact Name" value={load.carrier.contact.name} />
            <InfoRow label="Contact Phone" value={load.carrier.contact.phone} />
          </Grid2>
          <Grid2 size={6}>
            <InfoRow label="Equipment Type" value={load.carrier.equipment.type} />
            <InfoRow label="Equipment Length" value={`${load.carrier.equipment.length} ft`} />
          </Grid2>
        </Grid2>
      </DetailSection>

      <DetailSection title="Specifications">
        <Grid2 container spacing={2}>
          <Grid2 size={6}>
            <InfoRow label="Service Level" value={load.specifications.serviceLevel} />
            <InfoRow label="Temperature Min" value={`${load.specifications.temperature.min}°${load.specifications.temperature.unit}`} />
            <InfoRow label="Temperature Max" value={`${load.specifications.temperature.max}°${load.specifications.temperature.unit}`} />
          </Grid2>
          <Grid2 size={6}>
            <InfoRow label="Special Instructions" value={load.specifications.specialInstructions} />
          </Grid2>
        </Grid2>
      </DetailSection>
    </Box>
  );
};