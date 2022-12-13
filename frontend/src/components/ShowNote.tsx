import { Row, Col, Stack, Badge, Button, Container } from "react-bootstrap";
import { Link, useNavigate, useParams } from "react-router-dom";
import ReactMarkdown from 'react-markdown'
import { useEffect, useState } from "react";
import Note from "../types/Note";
import NoteService from "../service/NoteService";
import { OnDeleteNote } from "./OnDeleteNote";

const ShowNote: React.FC = () => {

    const params = useParams()

    const navigate = useNavigate()

    const [note, setNote] = useState<Note>({ id: 0, header: "", body: "", tags: [] })


    useEffect(() => {
        retrieveNote();
    }, []);

    const retrieveNote = () => {
        NoteService.get(`/api/v1/notes/${params.id}`)
            .then((response: any) => {
                setNote({
                    id: response.id,
                    header: response.header,
                    body: response.body,
                    tags: response.tags
                })
            })
            .catch((e: Error) => {
                console.log(e);
            });
    };

    return <Container className="my-4">
        <Row className="align-items-center mb-4">
            <Col>
                <h2>{note.header}</h2>
                {note.tags.length > 0 && (
                    <Stack gap={1} direction="horizontal" className="mb-2 flex-wrap">
                        {note.tags.map(tag => (
                            <Badge className="text-truncate" key={tag.id}>
                                {tag.label}
                            </Badge>
                        ))}
                    </Stack>
                )}
            </Col>
            <Col xs="auto">
                <Stack gap={2} direction="horizontal">
                    <Link to={`/${note.id}/edit`}>
                        <Button>edit</Button>
                    </Link>
                    <Link to="/">
                        <Button onClick={() => {
                            OnDeleteNote(note.id, note.tags)
                        }
                        }>delete</Button>
                    </Link>
                    <Link to="..">
                        <Button>back</Button>
                    </Link>
                </Stack>
            </Col>
        </Row>
        <ReactMarkdown>{note.body}</ReactMarkdown>
    </Container>
}

export default ShowNote;