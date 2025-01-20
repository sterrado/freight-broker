import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, CssBaseline, Box } from '@mui/material';
import { createTheme } from '@mui/material/styles';
import { Sidebar } from './components/Layout/Sidebar';
import { MainContent } from './components/Layout/MainContent';

import { LoadsTable } from './components/Loads/LoadsTable';
import { LoadDetail } from './components/Loads/LoadDetail';
import { CreateLoad } from './components/Loads/CreateLoad';
import { Settings } from './components/Settings/Settings';

const theme = createTheme();


function App() {

  return (
    <Router>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Box sx={{ display: 'flex' }}>
          <Sidebar />
          <MainContent>
            <Routes>
              <Route path="/" element={<LoadsTable />} />
              <Route path="/loads/:id" element={<LoadDetail />} />
              <Route path="/loads/create" element={<CreateLoad />} />
              <Route path="/settings" element={<Settings />} />
            </Routes>
          </MainContent>
        </Box>
      </ThemeProvider>
    </Router>
  );
}

export default App;