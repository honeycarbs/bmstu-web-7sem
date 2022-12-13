import { Dispatch, Key, SetStateAction, useEffect, useState } from "react"
import { Row, Col, Stack, Button, Form, Container, Modal, ToggleButton } from "react-bootstrap"
import { Link, NavigateFunction, useNavigate } from "react-router-dom"
import ReactSelect from "react-select"
import AccountService from "../service/AccountService"
import NoteService from "../service/NoteService"
import TagService from "../service/TagService"
import Note from "../types/Note"
import Tag from "../types/Tag"
import NoteCard from "./NoteCard"
import { OnUpdateNote } from "./OnUpdateNote"
import { OnUpdateTag } from "./OnUpdateTag"
import { CustomSelectStyle } from "./SelectStyles"

type EditTagsModalProps = {
    AvailableTags: Tag[],
    SetAvailableTags: Dispatch<Tag[]>
    show: boolean,
    handleClose: () => void
}


const NoteList: React.FC = () => {
    const [searchURL, setSearchURL] = useState<string>("?")
    const [selectedTags, setSelectedTags] = useState<Tag[]>([])
    const [AvailableTags, setAvailableTags] = useState<Tag[]>([])
    const TagList = AvailableTags
    const [editTagsOpen, setEditTagsOpen] = useState<boolean>(false)

    let navigate: NavigateFunction = useNavigate()

    const onLogout = () => {
        AccountService.logout()
        navigate("/login")
    }

    const [notes, setNotes] = useState<Array<Note>>([])

    useEffect(() => {
        FetchNotes()
        FetchTags()
    }, [searchURL, setNotes])

    useEffect(() => {
        window.location.reload
    }, [setNotes, setEditTagsOpen])

    const FetchNotes = () => {
        NoteService.getAll(searchURL)
            .then((response: any) => {
                setNotes(response.notes)
            })
            .catch((e: Error) => {
                console.log(e)
            })
    }

    const FetchTags = () => {
        TagService.getAll()
            .then((response: any) => {
                setAvailableTags(response.tags)
            })
            .catch((e: Error) => {
                console.log(e)
            })
    }

    return (
        <Container className="my-4">
            <Row className="align-items-center mb-4">
                <Col>
                    <h2>my neat.ly</h2>
                </Col>
                <Col xs="auto">
                    <Stack gap={2} direction="horizontal">
                        <Link to="/new">
                            <Button>add note</Button>
                        </Link>
                        <Link to="/">
                            <Button
                                onClick={() => setEditTagsOpen(true)}>
                                edit tags
                            </Button>
                        </Link>
                        <Link to="/register">
                            <Button onClick={onLogout}>log out</Button >
                        </Link>
                    </Stack>
                </Col>
            </Row>
            <Row>
                <Col>
                    <Form.Group className="mb-3" controlId="tags">
                        <ReactSelect
                            placeholder="tags"
                            styles={ CustomSelectStyle }
                            value={selectedTags.map(tag => {
                                return { label: tag.label, value: Number(tag.id) }
                            })}
                            onChange={tags => {
                                setSelectedTags(
                                    tags.map(tag => {
                                        return { label: tag.label, id: Number(tag.value), color: "BBAC9F" }
                                    })
                                )
                                setSearchURL(
                                    "?" +
                                    tags.map(tag => {
                                        return "tag=" + tag.label
                                    }).join("&")
                                )
                            }}
                            options={AvailableTags.map(tag => {
                                return { label: tag.label, value: tag.id }
                            })}
                            isMulti
                        />
                    </Form.Group>
                </Col>
            </Row>
            <Row xs={1} sm={2} lg={3} xl={4} className="g-3">
                {notes.map((note) => (
                    <Col key={note.id}>
                        <NoteCard id={note.id} header={note.header} tags={note.tags} />
                    </Col>
                ))}
            </Row>
            <EditTagsModal
                AvailableTags={TagList}
                SetAvailableTags={setAvailableTags}
                show={editTagsOpen} handleClose={() => {
                    setEditTagsOpen(false)
                    window.location.reload()
                }} />
        </Container>
    )
}

export default NoteList

function EditTagsModal({ AvailableTags, SetAvailableTags, show, handleClose }: EditTagsModalProps) {
    // const [AvailableTags, setAvailableTags] = useState<Tag[]>(AvailableTags)
    // console.log(AvailableTags)
    // console.log(AvailableTags)
    // const [updatedID, setUpdatedID] = useState(false)
    // useEffect() 
    // useEffect(() => {
    //     // FetchTags()
    //     // FetchNotes()
    // }, [SetAvailableTags])

    return <Modal show={show} onHide={handleClose}>
        <Modal.Header closeButton>
            <Modal.Title>Edit Tags</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            <Form>
                <Stack gap={2}>
                    {AvailableTags.map(tag => (
                        <Row key={tag.id}>
                            <Col>
                                <Form.Control
                                    type="text"
                                    value={tag.label}
                                    onChange={e => {
                                        SetAvailableTags(
                                            AvailableTags.map(t => {
                                                if (t.id == tag.id) {
                                                    const edited: Tag = { id: t.id, label: e.target.value }
                                                    OnUpdateTag(tag.id, edited)
                                                    return edited
                                                } else {
                                                    return { id: t.id, label: t.label }
                                                }
                                            })
                                        )
                                        // console.log("update time")
                                    }}
                                />
                            </Col>
                        </Row>
                    ))}
                </Stack>
            </Form>
        </Modal.Body>
    </Modal >
}