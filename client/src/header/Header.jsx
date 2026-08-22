import { useNavigate } from 'react-router-dom'
import Navbar from '../components/Navbar'
import Button from '../components/Button'

const Header = () => {

    const navigate = useNavigate()

  return (
    <header>
        <Navbar>
            <Button onClick={() => navigate('/login')}>Sign In</Button>
        </Navbar>
    </header>
  )
}

export default Header
