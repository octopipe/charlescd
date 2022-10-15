import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { styled, useTheme, Theme, CSSObject } from '@mui/material/styles';
import Box from '@mui/material/Box';
import MuiDrawer from '@mui/material/Drawer';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import MenuIcon from '@mui/icons-material/Menu';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { NavLink, Outlet } from 'react-router-dom';

const drawerWidth = 240;

const items = [
  {
    name: 'Home',
    to: '/',
    icon: 'home'
  },
  {
    name: 'Circles',
    to: '/circles',
    icon: ["far", "circle"]
  },
  {
    name: 'Modules',
    to: '/modules',
    icon: 'folder'
  },
  {
    name: 'Routes',
    to: '/routes',
    icon: 'signs-post'
  }
]


const openedMixin = (theme: Theme): CSSObject => ({
  width: drawerWidth,
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
  overflowX: 'hidden',
});

const closedMixin = (theme: Theme): CSSObject => ({
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  overflowX: 'hidden',
  width: `calc(${theme.spacing(9)} + 1px)`,
  [theme.breakpoints.up('sm')]: {
    width: `calc(${theme.spacing(10)} + 1px)`,
  },
});

const DrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-end',
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
}));


const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
  ({ theme, open }) => ({
    width: drawerWidth,
    flexShrink: 0,
    whiteSpace: 'nowrap',
    boxSizing: 'border-box',
    ...(open && {
      ...openedMixin(theme),
      '& .MuiDrawer-paper': openedMixin(theme),
    }),
    ...(!open && {
      ...closedMixin(theme),
      '& .MuiDrawer-paper': closedMixin(theme),
    }),
  }),
);

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  borderBottom: 'none',
  marginLeft: '80px',
  width: `calc(100% - 80px)`,
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Main = () => {
  const theme = useTheme();
  const [open, setOpen] = useState(false);

  const handleDrawerOpen = () => {
    setOpen(open => !open);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  return (
    <Box sx={{ display: 'flex' }}>
      
      <AppBar position='fixed' open={open} variant="outlined">
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            onClick={handleDrawerOpen}
            edge="start"
            
          >
            <MenuIcon />
          </IconButton>
          
        </Toolbar>
      </AppBar>
      <Drawer variant='permanent' open={open}>

        <Box sx={{p: 2}}>
          <a className='d-flex align-items-center mb-3 mb-md-0 me-md-auto text-white text-decoration-none'>
            <svg width="230px" height="38px" viewBox="0 0 275 38" version="1.1"><title>294EDC80-3BF7-4126-869D-3546802C6B47</title><defs><linearGradient x1="14.4090158%" y1="85.6829387%" x2="89.0132248%" y2="16.4103413%" id="linearGradient-1"><stop stopColor="#EF4123" offset="0%"></stop><stop stopColor="#F7941E" offset="100%"></stop></linearGradient><linearGradient x1="100%" y1="26.858854%" x2="21.1498798%" y2="83.2440491%" id="linearGradient-2"><stop stopColor="#8DC63F" offset="0%"></stop><stop stopColor="#40BA8D" offset="100%"></stop></linearGradient></defs><g id="Landing-page" stroke="none" strokeWidth="1" fill="none" fillRule="evenodd"><g id="LP---3.0---On-development" transform="translate(-269.000000, -151.000000)"><g id="banner"><g id="Brand/Charles/Color-White" transform="translate(269.000000, 151.000000)"><g id="Combined-Shape"><path d="M42.8003364,0.92967 L43.0922641,0.933337397 C49.2042059,1.08705766 54.1284439,6.05363354 54.1284439,12.1371367 L54.1284439,12.1371367 L54.1247368,12.4259687 C53.9693533,18.4731018 48.9490417,23.3452367 42.8003364,23.3452367 L42.8003364,23.3452367 L42.5083793,23.3415688 C36.3958268,23.1878294 31.4715888,18.2206398 31.4715888,12.1371367 C31.4715888,5.95707 36.5533925,0.92967 42.8003364,0.92967 L42.8003364,0.92967 Z M42.8003364,2.66120333 L42.5280157,2.66497806 C37.3712614,2.80809826 33.2212196,7.00284529 33.2212196,12.1371367 C33.2212196,17.3621367 37.5181542,21.6137033 42.8003364,21.6137033 L42.8003364,21.6137033 L43.0725926,21.6099286 C48.228153,21.4668082 52.3788131,17.2720505 52.3788131,12.1371367 L52.3788131,12.1371367 L52.3749975,11.8677951 C52.2303285,6.7674281 47.9901883,2.66120333 42.8003364,2.66120333 L42.8003364,2.66120333 Z" fill="#FFFFFF"></path><path d="M15.7325935,6.38013667 L16.0492462,6.38324168 C24.578849,6.55064876 31.4651869,13.4660322 31.4651869,21.94367 L31.4651869,21.94367 L31.462048,22.2569328 C31.2928141,30.6952211 24.3019716,37.5078367 15.7325935,37.5078367 L15.7325935,37.5078367 L15.4159407,37.5047314 C6.88633794,37.3373092 -3.51503357e-13,30.4213078 -3.51503357e-13,21.94367 C-3.51503357e-13,13.36137 7.05742056,6.38013667 15.7325935,6.38013667 L15.7325935,6.38013667 Z M15.7325935,9.44293667 L15.4260478,9.44656143 C8.59931269,9.60816051 3.09658411,15.1524646 3.09658411,21.94367 C3.09658411,28.83687 8.76479907,34.4450367 15.7325935,34.4450367 L15.7325935,34.4450367 L16.0391391,34.4414119 C22.8658742,34.2798127 28.3686028,28.7354994 28.3686028,21.94367 L28.3686028,21.94367 L28.3649392,21.6404332 C28.2016095,14.8873701 22.5979203,9.44293667 15.7325935,9.44293667 L15.7325935,9.44293667 Z" fill="url(#linearGradient-1)"></path><path d="M39.9322991,26.1173367 L40.1503582,26.1214122 C43.189201,26.2351846 45.6261215,28.7165661 45.6261215,31.7502033 L45.6261215,31.7502033 L45.6220018,31.965928 C45.5069983,34.9722382 42.9987646,37.38307 39.9322991,37.38307 L39.9322991,37.38307 L39.71424,37.3789944 C36.6753971,37.265222 34.2384766,34.7838405 34.2384766,31.7502033 C34.2384766,28.6443367 36.7928224,26.1173367 39.9322991,26.1173367 L39.9322991,26.1173367 Z M39.9322991,27.87737 L39.7311992,27.8824205 C37.6661918,27.9863745 36.0181963,29.6819544 36.0181963,31.7502033 C36.0181963,33.8858033 37.774229,35.6230367 39.9322991,35.6230367 L39.9322991,35.6230367 L40.1334553,35.617988 C42.1990071,35.5140706 43.8464019,33.8190658 43.8464019,31.7502033 L43.8464019,31.7502033 L43.8413004,31.5512545 C43.736296,29.5083171 42.0235496,27.87737 39.9322991,27.87737 L39.9322991,27.87737 Z" fill="url(#linearGradient-2)"></path></g></g></g></g></g></svg>
          </a>
        </Box>
        <Divider />
        <List sx={{p: 2}}>
          {items.map(item => (
            <NavLink to={item.to}>
              <ListItemButton
                sx={{
                  minHeight: 48,
                  justifyContent: open ? 'initial' : 'center',
                  px: 2.5,
                }}
              >
                <ListItemIcon
                  sx={{
                    minWidth: 0,
                    mr: open ? 3 : 'auto',
                    justifyContent: 'center',
                  }}
                >
                  <FontAwesomeIcon icon={item.icon as any} />
                </ListItemIcon>
                <ListItemText primary={item.name} sx={{ opacity: open ? 1 : 0 }} />
              </ListItemButton>
            </NavLink>
          ))}
        </List>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Outlet />
      </Box>
    </Box>
  )
}

export default Main 