import { Outlet, useNavigate, useParams } from "react-router-dom"
import NoteService from "../service/NoteService"

const NoteLayout: React.FC = () => {
    const params = useParams()
    const navigate = useNavigate()

    try {
        NoteService.get(`/api/v1/notes/${params.id}`).then(
            (response) => {
                console.log("ok")
            },
            (error) => {
                console.log("error")
                navigate("/")
            }
        )
    } catch (error) {
        console.log("ne ok")
        navigate("/")
    }

    return <Outlet />
}

export default NoteLayout