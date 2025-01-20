import React from 'react';
import { Box, styled } from '@mui/material';

const MainContentWrapper = styled(Box)(({ theme }) => ({
  flexGrow: 1,
  padding: theme.spacing(3),
  marginLeft: 240, // Same as drawer width
}));

export const MainContent: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <MainContentWrapper>{children}</MainContentWrapper>;
};