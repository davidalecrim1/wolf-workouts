declare module 'react-big-calendar' {
    import { ComponentType } from 'react'
    
    export interface CalendarProps {
        localizer: any
        events: any[]
        startAccessor: string
        endAccessor: string
        style?: object
    }

    export const momentLocalizer: (moment: any) => any
    export const Calendar: ComponentType<CalendarProps>
}