import React, { useState } from 'react'
import { Link, NavLink, useHistory, useLocation } from 'react-router-dom';
import { Button, Icon, Image, Menu } from 'semantic-ui-react';
//import logo from '../../../public/plus.png'

const Navbar = () => {
     const [isOpen, setOpen] = useState(false);
     const {pathname} = useLocation();
     const history = useHistory();
    
    return (
        <Menu secondary pointing>
            <Menu.Item as={Link} to="/" style={{fontSize:24}}>All books</Menu.Item>
           {pathname==="/" && (
            <Menu.Item position="right">
                <Button as={Link} to="/book/create" icon  primary basic>
                    <Icon name="add"></Icon>
                    New Book
                </Button>
            </Menu.Item>
            )}
            
            <Menu.Item as={Link} to="/users">
                    <Button icon  color="yellow" basic>
                        <Icon name="user"></Icon>
                        Users
                    </Button>
            </Menu.Item>
            <Menu.Item position="left">
                <Button as={Link} to="/users/create" icon   basic color="green">
                    <Icon name="add"></Icon>
                    New User
                </Button>
            </Menu.Item>
            {pathname==="/" && (
                <Menu.Item>
                    <Button icon  color="red" basic>
                        <Icon name="log out"></Icon>
                        logout
                    </Button>
                </Menu.Item>
            )}
        </Menu>
    )
}

export default Navbar