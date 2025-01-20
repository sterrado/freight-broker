import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  Box, 
  Button,
  Paper,
  Typography 
} from '@mui/material';
import { 
  DataGrid, 
  GridColDef, 
  GridPaginationModel,
  GridRowParams,
} from '@mui/x-data-grid';
import { Add as AddIcon } from '@mui/icons-material';
import { loadService } from '../../services/api';
import { Load } from '../../types/load.types';

const columns: GridColDef<Load>[] = [
    { 
      field: 'id',  // matches the Load interface
      headerName: 'Freight ID', 
      width: 130,
      valueGetter: (params: { row: Load }) => params.row.freightLoadID
    },
    { 
      field: 'customer',  
      headerName: 'Customer', 
      width: 180,
      valueGetter: (params: { row: Load }) => params.row.customer.name
    },
    {
      field: 'status',
      headerName: 'Status',
      width: 130,
      valueGetter: (params: { row: Load }) => params.row.status.code.value
    },
    {
      field: 'pickup',
      headerName: 'Pickup Location',
      width: 150,
      valueGetter: (params: { row: Load }) => params.row.pickup.address.city
    },
    {
      field: 'consignee',  // changed from delivery to match interface
      headerName: 'Delivery Location',
      width: 150,
      valueGetter: (params: { row: Load }) => params.row.consignee.address.city
    },
    {
      field: 'rateData',
      headerName: 'Total Rate',
      width: 130,
      type: 'number',
      valueGetter: (params: { row: Load }) => params.row.rateData.totalRate,
      renderCell: (params: { row: Load }) => `$${params.row.rateData.totalRate.toLocaleString('en-US', { minimumFractionDigits: 2 })}`
    },
    {
      field: 'createdAt',
      headerName: 'Created',
      width: 180,
      valueGetter: (params: { row: Load }) => new Date(params.row.createdAt).toLocaleString()
    }
  ];

export const LoadsTable = () => {
  const navigate = useNavigate();
  const [loads, setLoads] = useState<Load[]>([]);
  const [loading, setLoading] = useState(true);
  const [totalRows, setTotalRows] = useState(0);
  const [paginationModel, setPaginationModel] = useState<GridPaginationModel>({
    pageSize: 10,
    page: 0,
  });

  useEffect(() => {
    const fetchLoads = async () => {
      try {
        setLoading(true);
        const response = await loadService.getLoads(
          paginationModel.page + 1,
          paginationModel.pageSize
        );
        setLoads(response.loads);
        setTotalRows(response.total);
      } catch (error) {
        console.error('Error fetching loads:', error);
        // TODO: Add error handling/notification
      } finally {
        setLoading(false);
      }
    };

    fetchLoads();
  }, [paginationModel]);

  const handleCreateLoad = () => {
    navigate('/loads/create');
  };

  const handleRowClick = (params: GridRowParams<Load>) => {
    navigate(`/loads/${params.row.id}`);
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h5" component="h1">
          Loads
        </Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={handleCreateLoad}
        >
          Create Load
        </Button>
      </Box>
      <Paper sx={{ height: 'calc(100vh - 200px)', width: '100%' }}>
        <DataGrid
          rows={loads}
          columns={columns}
          loading={loading}
          rowCount={totalRows}
          pageSizeOptions={[5, 10, 25, 50]}
          paginationMode="server"
          paginationModel={paginationModel}
          onPaginationModelChange={setPaginationModel}
          onRowClick={handleRowClick}
          disableRowSelectionOnClick
          autoHeight
        />
      </Paper>
    </Box>
  );
};